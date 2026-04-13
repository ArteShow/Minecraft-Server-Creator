import React, { useEffect, useMemo, useRef, useState } from "react";
import { createPortal } from "react-dom";
import {
  Boxes,
  Check,
  ChevronLeft,
  ChevronRight,
  Copy,
  CreditCard,
  Crown,
  Database,
  Download,
  HardDrive,
  KeyRound,
  LayoutDashboard,
  Lock,
  LogIn,
  Menu,
  MessageSquare,
  Play,
  Plus,
  RefreshCw,
  Search,
  Server,
  Settings,
  Shield,
  Square,
  TerminalSquare,
  Trash2,
  UserPlus,
  Wifi,
  X,
} from "lucide-react";

import logoFull from "../nginx/html/logo-full.png";
import logoIcon from "../nginx/html/logo-icon.png";

const API_BASE = "/api/v1";
const SHOW_PLUGIN_CATALOG = false;
const FIRST_SERVER_PORT = 25565;

const plans = [
  { name: "Starter", ram: 2, cores: 2, backups: 0, storage: 5, price: 5.99, featured: false },
  { name: "Standard", ram: 4, cores: 3, backups: 1, storage: 15, price: 8.99, featured: true },
  { name: "Premium", ram: 5, cores: 4, backups: 2, storage: 20, price: 10.99, featured: false },
];

const addOns = [
  { key: "ram1", label: "+1 GB RAM", price: 1.99, ram: 1, storage: 0 },
  { key: "ssd5", label: "+5 GB SSD", price: 1.49, ram: 0, storage: 5 },
  { key: "ssd10", label: "+10 GB SSD", price: 2.0, ram: 0, storage: 10 },
];

const versions = ["1.8.9", "1.12.2", "1.16.5", "1.20.6", "1.21.1", "1.21.2"];
const softwareOptions = ["Vanilla", "Spigot", "Paper"];
const pluginCapableServerTypes = new Set(["Spigot", "Paper"]);
const PLUGIN_PAGE_SIZE = 12;

const fallbackMods = [
  {
    id: "sodium",
    name: "Sodium",
    source: "Modrinth",
    type: "Performance",
    description: "Performance-focused rendering optimization mod.",
  },
  {
    id: "lithium",
    name: "Lithium",
    source: "Modrinth",
    type: "Optimization",
    description: "General server and game logic optimizations.",
  },
  {
    id: "fabric-api",
    name: "Fabric API",
    source: "Modrinth",
    type: "Core",
    description: "Core hooks and shared API for Fabric mods.",
  },
];

function cn(...classes) {
  return classes.filter(Boolean).join(" ");
}

function money(value) {
  return new Intl.NumberFormat("en-EN", { style: "currency", currency: "EUR" }).format(value);
}

function safeJsonParse(value, fallback) {
  try {
    return value ? JSON.parse(value) : fallback;
  } catch {
    return fallback;
  }
}

function extractErrorMessage(rawText, fallbackMessage) {
  const trimmed = String(rawText || "").trim();
  if (!trimmed) return fallbackMessage;

  if (trimmed.startsWith("<")) {
    const titleMatch = trimmed.match(/<title>(.*?)<\/title>/i);
    const headingMatch = trimmed.match(/<h1>(.*?)<\/h1>/i);
    return titleMatch?.[1] || headingMatch?.[1] || fallbackMessage;
  }

  const parsed = safeJsonParse(trimmed, null);
  if (parsed && typeof parsed === "object") {
    return parsed.message || parsed.error || parsed.raw || trimmed;
  }

  return trimmed;
}

function normalizePluginToken(value) {
  return String(value || "")
    .toLowerCase()
    .replace(/\.jar$/i, "")
    .replace(/[^a-z0-9]/g, "");
}

function isPluginAlreadyInstalled(plugin, installedFiles) {
  const nameToken = normalizePluginToken(plugin?.name);
  const slugToken = normalizePluginToken(plugin?.slug);

  return (installedFiles || []).some((fileName) => {
    const fileToken = normalizePluginToken(fileName);
    return (nameToken && fileToken.includes(nameToken)) || (slugToken && fileToken.includes(slugToken));
  });
}

async function searchModrinthPlugins(query, page = 0) {
  const url = new URL("https://api.modrinth.com/v2/search");
  url.searchParams.set("query", (query || "").trim());
  url.searchParams.set("limit", String(PLUGIN_PAGE_SIZE));
  url.searchParams.set("offset", String(Math.max(page, 0) * PLUGIN_PAGE_SIZE));
  url.searchParams.set("index", "downloads");
  url.searchParams.set("facets", JSON.stringify([["project_type:plugin"]]));

  const response = await fetch(url.toString());
  if (!response.ok) {
    throw new Error(`Modrinth search failed with status ${response.status}`);
  }

  const data = await response.json();
  const items = (data.hits || []).map((item) => ({
    id: item.project_id,
    name: item.title,
    description: item.description,
    author: item.author,
    downloads: item.downloads,
    slug: item.slug,
    iconUrl: item.icon_url || "",
  }));

  return {
    items,
    hasMore: items.length === PLUGIN_PAGE_SIZE,
  };
}

async function fetchLatestPluginFile(projectId, serverVersion) {
  const url = new URL(`https://api.modrinth.com/v2/project/${projectId}/version`);
  url.searchParams.set("loaders", JSON.stringify(["paper", "spigot", "bukkit", "purpur"]));
  if (serverVersion) {
    url.searchParams.set("game_versions", JSON.stringify([serverVersion]));
  }

  const response = await fetch(url.toString());
  if (!response.ok) {
    throw new Error(`Modrinth versions failed with status ${response.status}`);
  }

  const versionsData = await response.json();
  const hasFiles = (entry) => Array.isArray(entry.files) && entry.files.length > 0;
  const releaseVersion = versionsData.find((entry) => entry.version_type === "release" && hasFiles(entry));
  const betaVersion = versionsData.find((entry) => entry.version_type === "beta" && hasFiles(entry));
  const alphaVersion = versionsData.find((entry) => entry.version_type === "alpha" && hasFiles(entry));
  const version = releaseVersion || betaVersion || alphaVersion || versionsData.find(hasFiles);
  if (!version) {
    throw new Error("No compatible plugin file found for this server version.");
  }

  const primaryFile = version.files.find((file) => file.primary) || version.files[0];
  if (!primaryFile?.url) {
    throw new Error("Plugin version has no downloadable file.");
  }

  const fileResponse = await fetch(primaryFile.url);
  if (!fileResponse.ok) {
    throw new Error(`Plugin download failed with status ${fileResponse.status}`);
  }

  const blob = await fileResponse.blob();
  return new File([blob], primaryFile.filename || `${projectId}.jar`, { type: blob.type || "application/java-archive" });
}

function decodeJwt(token) {
  try {
    const payload = token.split(".")[1];
    if (!payload) return null;
    const normalized = payload.replace(/-/g, "+").replace(/_/g, "/");
    const json = decodeURIComponent(
      atob(normalized)
        .split("")
        .map((char) => `%${(`00${char.charCodeAt(0).toString(16)}`).slice(-2)}`)
        .join(""),
    );
    return JSON.parse(json);
  } catch {
    return null;
  }
}

async function apiFetch(path, { method = "GET", body, token } = {}) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (token) headers.Authorization = `Bearer ${token}`;

  const response = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  const text = await response.text();
  const data = text ? safeJsonParse(text, { raw: text }) : {};

  if (response.status === 401) {
    window.dispatchEvent(
      new CustomEvent("easy2host:unauthorized", {
        detail: { path, method, message: data?.message || data?.error || "Unauthorized" },
      }),
    );
  }

  if (!response.ok) {
    throw new Error(data?.message || data?.error || data?.raw || `Request failed with status ${response.status}`);
  }

  return data;
}

async function apiFetchBlob(path, { method = "GET", body, token } = {}) {
  const headers = {};
  if (body !== undefined) headers["Content-Type"] = "application/json";
  if (token) headers.Authorization = `Bearer ${token}`;

  const response = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    const text = await response.text();
    const data = text ? safeJsonParse(text, { raw: text }) : {};
    if (response.status === 401) {
      window.dispatchEvent(
        new CustomEvent("easy2host:unauthorized", {
          detail: { path, method, message: data?.message || data?.error || "Unauthorized" },
        }),
      );
    }
    throw new Error(data?.message || data?.error || data?.raw || `Request failed with status ${response.status}`);
  }

  return response.blob();
}

function popClass() {
  return "transition-all duration-200 motion-safe:hover:-translate-y-0.5 sm:motion-safe:hover:-translate-y-1 sm:motion-safe:hover:scale-[1.01]";
}

function clampPercent(value) {
  if (!Number.isFinite(value)) return 0;
  if (value < 0) return 0;
  if (value > 100) return 100;
  return value;
}

function makeTicketId() {
  return `TCK-${Math.random().toString(36).slice(2, 8).toUpperCase()}`;
}

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(url);
}

function buildConsoleWsUrl(serverId, token) {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}${API_BASE}/server/console/ws?server_id=${encodeURIComponent(serverId)}&token=${encodeURIComponent(token)}`;
}

function GlassCard({ className = "", children }) {
  return (
    <div
      className={cn(
        "min-w-0 overflow-hidden rounded-[2rem] border border-white/10 bg-white/5 backdrop-blur-xl shadow-2xl shadow-black/20",
        className,
      )}
    >
      {children}
    </div>
  );
}

function Brand({ compact = false }) {
  return (
    <div className="flex items-center gap-3">
      <img src={compact ? logoIcon : logoFull} alt="easy2host" className={compact ? "h-11 w-11 rounded-2xl" : "h-10 w-auto sm:h-11"} />
      {!compact && (
        <div>
          <div className="display-font text-lg font-bold tracking-tight text-white">easy2host</div>
          <div className="text-[11px] uppercase tracking-[0.32em] text-slate-400">Minecraft hosting made easy</div>
        </div>
      )}
    </div>
  );
}

function SectionIntro({ eyebrow, title, text }) {
  return (
    <div className="max-w-3xl">
      <p className="mb-3 text-sm font-semibold uppercase tracking-[0.24em] text-cyan-300/80">{eyebrow}</p>
      <h2 className="display-font text-3xl font-bold tracking-tight text-white sm:text-4xl">{title}</h2>
      <p className="mt-4 text-base leading-7 text-slate-300">{text}</p>
    </div>
  );
}

function Modal({ open, onClose, title, subtitle, children }) {
  // Trap focus and stop background scroll while open
  useEffect(() => {
    if (!open) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => { document.body.style.overflow = prev; };
  }, [open]);

  if (!open) return null;

  return createPortal(
    <div
      className="fixed inset-0 z-[120] flex items-end justify-center bg-black/75 p-2 sm:items-center sm:p-4"
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <GlassCard
        className="relative z-10 max-h-[92vh] w-full max-w-3xl overflow-auto rounded-t-[2.2rem] p-5 sm:rounded-[2rem] sm:p-8"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mb-6 flex items-start justify-between gap-4">
          <div>
            <h2 className="display-font text-2xl font-bold text-white">{title}</h2>
            <p className="mt-2 text-sm text-slate-300 sm:text-base">{subtitle}</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="shrink-0 rounded-2xl border border-white/10 bg-white/5 p-2 text-white transition hover:bg-white/10"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        {children}
      </GlassCard>
    </div>,
    document.body,
  );
}

function Sidebar({ items, active, setActive, title, subtitle }) {
  return (
    <GlassCard className="p-4 lg:sticky lg:top-6 lg:h-fit">
      <div className="mb-5 px-2">
        <div className="text-xs uppercase tracking-[0.24em] text-cyan-300/80">{subtitle}</div>
        <div className="display-font mt-2 text-xl font-bold text-white">{title}</div>
      </div>
      <div className="space-y-2">
        {items.map((item) => {
          const Icon = item.icon;
          return (
            <button
              key={item.key}
              onClick={() => setActive(item.key)}
              className={cn(
                "flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-left",
                popClass(),
                active === item.key ? "bg-cyan-400/12 text-white" : "text-slate-300 hover:bg-white/5",
              )}
            >
              <Icon className="h-5 w-5" />
              <span>{item.label}</span>
            </button>
          );
        })}
      </div>
    </GlassCard>
  );
}

function ToastStack({ notices, setNotices }) {
  useEffect(() => {
    if (notices.length === 0) return;
    const timer = setTimeout(() => {
      setNotices((prev) => prev.slice(0, -1));
    }, 5000);
    return () => clearTimeout(timer);
  }, [notices, setNotices]);

  if (notices.length === 0) return null;

  return createPortal(
    <div className="fixed bottom-5 right-5 z-[200] flex flex-col gap-2 max-w-sm w-full">
      {notices.map((notice) => (
        <div
          key={notice.id}
          className={cn(
            "flex items-start justify-between gap-3 rounded-2xl border px-4 py-3 text-sm shadow-2xl backdrop-blur-xl",
            notice.type === "error"
              ? "border-rose-400/20 bg-slate-950/95 text-rose-200"
              : notice.type === "info"
              ? "border-sky-400/20 bg-slate-950/95 text-sky-200"
              : "border-emerald-400/20 bg-slate-950/95 text-emerald-200",
          )}
        >
          <span className="flex-1">{notice.text}</span>
          <button
            type="button"
            onClick={() => setNotices((prev) => prev.filter((n) => n.id !== notice.id))}
            className="shrink-0 text-slate-400 hover:text-white"
          >
            <X className="h-4 w-4" />
          </button>
        </div>
      ))}
    </div>,
    document.body,
  );
}

function PurchaseFlow({ open, plan, onClose, currentUser, onRequireAuth, onComplete }) {
  const [step, setStep] = useState(1);
  const [selectedAddOns, setSelectedAddOns] = useState([]);
  const [setup, setSetup] = useState({ name: "", version: versions.at(-1), software: "Paper" });
  const [billing, setBilling] = useState({ fullName: currentUser?.displayName || "", country: "Austria" });
  const [creating, setCreating] = useState(false);
  const [createStep, setCreateStep] = useState("");

  useEffect(() => {
    if (open) {
      setBilling({ fullName: currentUser?.displayName || "", country: "Austria" });
      setStep(1);
      setSelectedAddOns([]);
      setSetup({ name: "", version: versions.at(-1), software: "Paper" });
      setCreating(false);
      setCreateStep("");
    }
  }, [open, currentUser]);

  if (!open || !plan) return null;

  const selected = addOns.filter((addon) => selectedAddOns.includes(addon.key));
  const total = plan.price + selected.reduce((sum, addon) => sum + addon.price, 0);

  const next = () => {
    if (!currentUser) {
      onRequireAuth();
      onClose();
      return;
    }
    if (step === 2 && !billing.fullName.trim()) return;
    if (step === 3 && (!setup.name.trim() || !setup.version.trim())) return;
    if (step < 3) {
      setStep((value) => value + 1);
      return;
    }

    setCreating(true);
    onClose();
    onComplete({ plan, addons: selected, setup, billing }, setCreateStep).finally(() => setCreating(false));
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={`Buy ${plan.name}`}
      subtitle="Customise your plan, review the order, then launch your server."
    >
      <div className="mb-6 flex flex-wrap gap-3">
        {["Extras", "Billing", "Server Setup"].map((label, index) => (
          <button
            key={label}
            type="button"
            onClick={() => !creating && setStep(index + 1)}
            className={cn(
              "rounded-full px-4 py-2 text-sm",
              creating ? "cursor-default" : "cursor-pointer",
              step === index + 1 ? "bg-cyan-400/15 text-cyan-300" : "bg-white/5 text-slate-400",
            )}
          >
            {label}
          </button>
        ))}
      </div>

      {step === 1 && (
        <div className="space-y-6">
          <GlassCard className="p-5">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <div className="display-font text-xl font-bold text-white">{plan.name}</div>
                <div className="mt-2 grid gap-1 text-sm text-slate-300 sm:grid-cols-2">
                  <div>{plan.ram} GB RAM</div>
                  <div>{plan.storage} GB SSD</div>
                  <div>{plan.cores} Cores</div>
                  <div>{plan.backups} included backups</div>
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm text-slate-400">Base price</div>
                <div className="display-font text-2xl font-bold text-white">{money(plan.price)}</div>
                <div className="mt-2 text-xs text-slate-400">
                  {selected.length > 0 ? `+ ${selected.length} add-on${selected.length > 1 ? "s" : ""}` : "No add-ons selected"}
                </div>
                <div className="display-font mt-2 text-lg font-bold text-cyan-300">Total: {money(total)}</div>
              </div>
            </div>
          </GlassCard>
          <div className="grid gap-4 sm:grid-cols-3">
            {addOns.map((addon) => {
              const active = selectedAddOns.includes(addon.key);
              return (
                <button
                  key={addon.key}
                  onClick={() =>
                    setSelectedAddOns((previous) =>
                      previous.includes(addon.key)
                        ? previous.filter((entry) => entry !== addon.key)
                        : [...previous, addon.key],
                    )
                  }
                  className={cn(
                    "rounded-3xl border p-5 text-left",
                    popClass(),
                    active ? "border-cyan-300/30 bg-cyan-400/10" : "border-white/10 bg-white/5",
                  )}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="font-semibold text-white">{addon.label}</div>
                      <div className="mt-1 text-sm text-slate-400">{money(addon.price)} / month</div>
                    </div>
                    <div
                      className={cn(
                        "flex h-7 w-7 items-center justify-center rounded-full border",
                        active ? "border-cyan-300 bg-cyan-300/15 text-cyan-300" : "border-white/10 text-slate-500",
                      )}
                    >
                      <Check className="h-4 w-4" />
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {step === 2 && (
        <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
          <div className="space-y-4">
            <input
              value={billing.fullName}
              onChange={(event) => setBilling((current) => ({ ...current, fullName: event.target.value }))}
              placeholder="Full name"
              className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
            />
            <input
              value={billing.country}
              onChange={(event) => setBilling((current) => ({ ...current, country: event.target.value }))}
              placeholder="Country"
              className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
            />
            <div className="rounded-2xl border border-cyan-300/20 bg-cyan-400/10 px-4 py-3 text-sm text-cyan-100">
              Review your order on the right. The server will be provisioned after you complete the final step.
            </div>
          </div>
          <GlassCard className="p-5">
            <div className="mb-4 flex items-center gap-2 text-white">
              <CreditCard className="h-5 w-5 text-cyan-300" />
              Order Summary
            </div>
            <div className="space-y-3 text-sm text-slate-300">
              <div className="flex items-center justify-between">
                <span>{plan.name}</span>
                <span>{money(plan.price)}</span>
              </div>
              {selected.map((addon) => (
                <div key={addon.key} className="flex items-center justify-between">
                  <span>{addon.label}</span>
                  <span>{money(addon.price)}</span>
                </div>
              ))}
              <div className="flex items-center justify-between border-t border-white/10 pt-3 text-base font-semibold text-white">
                <span>Total</span>
                <span>{money(total)}</span>
              </div>
            </div>
          </GlassCard>
        </div>
      )}

      {step === 3 && (
        <div className="grid gap-4">
          <input
            value={setup.name}
            onChange={(event) => setSetup((current) => ({ ...current, name: event.target.value }))}
            placeholder="Server name"
            className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
          />
          <select
            value={setup.software}
            onChange={(event) => setSetup((current) => ({ ...current, software: event.target.value }))}
            className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none"
          >
            {softwareOptions.map((software) => (
              <option key={software}>{software}</option>
            ))}
          </select>
          <div className="grid gap-2">
            <input
              value={setup.version}
              onChange={(event) => setSetup((current) => ({ ...current, version: event.target.value.trim() }))}
              list="minecraft-version-list"
              placeholder="Minecraft version (e.g. 1.21.11)"
              className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
            />
            <datalist id="minecraft-version-list">
              {versions.map((version) => (
                <option key={version} value={version} />
              ))}
            </datalist>
            <div className="text-xs text-slate-400">You can type any version manually, for example 1.21.11.</div>
          </div>
            <div className="rounded-2xl border border-dashed border-cyan-300/25 bg-cyan-400/5 px-4 py-3 text-slate-300">
              Pick the Minecraft version and server software for your new server. Plugin installs are available for Spigot and Paper servers.
            </div>
        </div>
      )}

      <div className="mt-8 flex flex-wrap items-center justify-between gap-4">
        <div className="min-w-0 flex-1">
          {creating ? (
            <div className="rounded-2xl border border-cyan-300/20 bg-cyan-400/10 px-4 py-3 text-sm text-cyan-100">
              <div className="font-semibold">Creating your server…</div>
              {createStep && <div className="mt-1 text-slate-300">{createStep}</div>}
            </div>
          ) : (
            <div className="text-sm text-slate-400">
              {currentUser ? `Signed in as ${currentUser.username}` : "Create an account or sign in to continue."}
            </div>
          )}
        </div>
        <div className="flex gap-3">
          {step > 1 && !creating && (
            <button
              onClick={() => setStep((value) => value - 1)}
              className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 font-semibold text-white"
            >
              Back
            </button>
          )}
          <button
            disabled={creating}
            onClick={next}
            className={cn(
              "rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60",
              !creating && popClass(),
            )}
          >
            {creating ? "Creating…" : !currentUser ? "Sign in to continue" : step === 3 ? "Create Server" : "Continue"}
          </button>
        </div>
      </div>
    </Modal>
  );
}

function AuthScreen({ mode, busy, error, onSubmit, setScreen }) {
  const signup = mode === "signup";
  const [form, setForm] = useState({ displayName: "", username: "", email: "", password: "" });

  return (
    <div className="mx-auto flex min-h-[calc(100vh-120px)] max-w-7xl items-center px-4 py-8 sm:px-6 sm:py-12 lg:px-16">
      <div className="grid w-full gap-6 lg:grid-cols-[0.95fr_1.05fr] lg:gap-8">
        <GlassCard className="p-5 sm:p-8">
          <div className="mb-5 flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">
            {signup ? <UserPlus className="h-7 w-7" /> : <LogIn className="h-7 w-7" />}
          </div>
          <h1 className="display-font text-2xl font-bold text-white sm:text-3xl">
            {signup ? "Create your account" : "Sign in"}
          </h1>
          <p className="mt-3 text-slate-300">
            {signup
              ? "Create a customer account with username, email, and password."
              : "Sign in with your username and password. Admin credentials open the admin dashboard automatically."}
          </p>
          <div className="mt-8 space-y-4">
            {signup && (
              <input
                value={form.displayName}
                onChange={(event) => setForm((current) => ({ ...current, displayName: event.target.value }))}
                placeholder="Display name"
                className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
              />
            )}
            <input
              value={form.username}
              onChange={(event) => setForm((current) => ({ ...current, username: event.target.value }))}
              placeholder="Username"
              className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
            />
            {signup && (
              <input
                value={form.email}
                onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
                placeholder="Email"
                className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
              />
            )}
            <input
              value={form.password}
              onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
              type="password"
              placeholder="Password"
              className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
            />
            {error && <div className="rounded-2xl border border-rose-400/20 bg-rose-400/10 px-4 py-3 text-sm text-rose-200">{error}</div>}
            <button
              disabled={busy}
              onClick={() => onSubmit(form)}
              className={cn(
                "w-full rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60",
                popClass(),
              )}
            >
              {busy ? "Please wait..." : signup ? "Create Account" : "Sign In"}
            </button>
          </div>
          <div className="mt-6 text-sm text-slate-400">
            {signup ? "Already have an account?" : "Need an account?"}{" "}
            <button onClick={() => setScreen(signup ? "signin" : "signup")} className="font-semibold text-cyan-300">
              {signup ? "Sign in" : "Create one"}
            </button>
          </div>
          {!signup && (
            <div className="mt-4 space-y-3">
              <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4 text-sm text-slate-300">
                <div className="font-semibold text-white">Quick links</div>
                <div className="mt-2 flex flex-wrap gap-3">
                  <a href="/dashboard.html" className="text-cyan-300 hover:text-cyan-200">Server Dashboard</a>
                  <a href="/admin.html" className="text-cyan-300 hover:text-cyan-200">Admin Dashboard</a>
                </div>
              </div>
            </div>
          )}
        </GlassCard>
        <GlassCard className="hidden overflow-hidden lg:block">
          <div className="h-full bg-[radial-gradient(circle_at_top_left,rgba(103,232,249,0.16),transparent_38%),linear-gradient(160deg,rgba(8,17,31,0.94),rgba(10,20,34,0.92))] p-8">
            <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-cyan-400/20 bg-cyan-400/10 px-4 py-2 text-sm text-cyan-200">
              easy2host — Simple Minecraft hosting
            </div>
            <div className="space-y-4">
              {[
                [Crown, "Simple, fast setup", "Create an account, pick a plan, and your server is ready in minutes."],
                [Server, "Full server control", "Start, stop, and manage your Minecraft server from one clean dashboard."],
                [HardDrive, "Backup anytime", "Create and download backups directly from your dashboard whenever you need."],
                [Shield, "Your servers stay yours", "Each account only sees and controls its own servers. Nothing shared."],
                [Settings, "Your choice of software", "Pick Vanilla, Spigot, or Paper — and any supported Minecraft version."],
              ].map(([Icon, title, text]) => (
                <div key={title} className="flex items-start gap-3 rounded-2xl border border-white/10 bg-white/5 p-4">
                  <Icon className="mt-0.5 h-5 w-5 shrink-0 text-cyan-300" />
                  <div>
                    <div className="font-semibold text-white">{title}</div>
                    <div className="mt-1 text-sm text-slate-300">{text}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </GlassCard>
      </div>
    </div>
  );
}

function PluginCatalogPanel({ search, setSearch }) {
  const filteredMods = fallbackMods.filter((item) => `${item.name} ${item.source} ${item.type}`.toLowerCase().includes(search.toLowerCase()));

  return (
    <GlassCard className="p-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Catalog</div>
          <h2 className="display-font mt-2 text-2xl font-bold text-white">Plugins & Mods</h2>
          <p className="mt-2 text-slate-400">Coming soon. Browse and install popular plugins and mods to enhance your server.</p>
        </div>
        <div className="relative w-full max-w-sm">
          <Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" />
          <input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            className="w-full rounded-2xl border border-white/10 bg-white/5 py-3 pl-11 pr-4 text-white outline-none"
            placeholder="Search plugins or mods"
          />
        </div>
      </div>
      <div className="mt-8 grid gap-4 xl:grid-cols-3">
        {filteredMods.map((item) => (
          <div key={item.id} className="rounded-2xl border border-white/10 bg-white/5 p-4 text-slate-300">
            <div className="font-semibold text-white">{item.name}</div>
            <div className="mt-1 text-sm text-slate-400">{item.source} · {item.type}</div>
            <p className="mt-3 text-sm leading-6">{item.description}</p>
          </div>
        ))}
      </div>
    </GlassCard>
  );
}

function CustomerDashboard({ currentUser, token, servers, setServers, notices, setNotices, apiHealthy, logout, tickets, setTickets, onBuyServer }) {
  const [active, setActive] = useState("servers");
  const [selectedServerId, setSelectedServerId] = useState(null);
  const [search, setSearch] = useState("");
  const [serverBusy, setServerBusy] = useState({});
  const [backupBusy, setBackupBusy] = useState({});
  const [serversRefreshing, setServersRefreshing] = useState(false);
  const [backupListMap, setBackupListMap] = useState({});
  const [selectedBackupMap, setSelectedBackupMap] = useState({});
  const [consoleTextMap, setConsoleTextMap] = useState({});
  const [playerCountMap, setPlayerCountMap] = useState({});
  const [consoleStateMap, setConsoleStateMap] = useState({});
  const [consoleWindowOpen, setConsoleWindowOpen] = useState(false);
  const [consoleCommand, setConsoleCommand] = useState("");
  const [consoleSending, setConsoleSending] = useState(false);
  const [consoleRefreshTick, setConsoleRefreshTick] = useState(0);
  const [uploadBackupModalOpen, setUploadBackupModalOpen] = useState(false);
  const [uploadWorldModalOpen, setUploadWorldModalOpen] = useState(false);
  const [pluginSearch, setPluginSearch] = useState("");
  const [pluginPage, setPluginPage] = useState(0);
  const [pluginHasMore, setPluginHasMore] = useState(false);
  const [pluginResults, setPluginResults] = useState([]);
  const [pluginLoading, setPluginLoading] = useState(false);
  const [pluginInstalling, setPluginInstalling] = useState({});
  const [installedPluginsMap, setInstalledPluginsMap] = useState({});
  const [installedPluginsLoading, setInstalledPluginsLoading] = useState({});
  const [pluginDeleting, setPluginDeleting] = useState({});
  const [powerUsageMap, setPowerUsageMap] = useState({});
  const [powerUsageLoading, setPowerUsageLoading] = useState({});
  const uploadBackupFileRef = useRef(null);
  const uploadWorldFileRef = useRef(null);
  const consoleOutputRef = useRef(null);

  const ownedServers = useMemo(
    () => servers.filter((server) => server.ownerId === currentUser.ownerKey),
    [servers, currentUser.ownerKey],
  );
  const selectedServer = ownedServers.find((server) => server.server_id === selectedServerId) || ownedServers[0] || null;

  useEffect(() => {
    if (!selectedServerId && ownedServers[0]) {
      setSelectedServerId(ownedServers[0].server_id);
    }
    if (selectedServerId && !ownedServers.some((server) => server.server_id === selectedServerId)) {
      setSelectedServerId(ownedServers[0]?.server_id ?? null);
    }
  }, [ownedServers, selectedServerId]);

  const setBusy = (serverId, value) => setServerBusy((previous) => ({ ...previous, [serverId]: value }));
  const setBackupLoading = (serverId, value) => setBackupBusy((previous) => ({ ...previous, [serverId]: value }));
  const pushNotice = (type, text) => setNotices((previous) => [{ id: Date.now() + Math.random(), type, text }, ...previous].slice(0, 5));

  const updateServer = (serverId, updater) => {
    setServers((previous) => previous.map((server) => (server.server_id === serverId ? { ...server, ...updater(server) } : server)));
  };

  const refreshOwnedServers = async ({ silent = false } = {}) => {
    if (!token || !currentUser?.ownerKey) return;

    const trackedServers = servers.filter((server) => server.ownerId === currentUser.ownerKey);
    if (trackedServers.length === 0) {
      if (!silent) pushNotice("success", "No tracked servers to refresh.");
      return;
    }

    setServersRefreshing(true);
    try {
      const results = await Promise.allSettled(
        trackedServers.map(async (server) => {
          const response = await apiFetch("/server/getStats", {
            method: "POST",
            body: { server_id: server.server_id, key: "Online" },
            token,
          });

          return {
            serverId: server.server_id,
            status: String(response?.value).toLowerCase() === "true" ? "online" : "offline",
          };
        }),
      );

      const nextStatusMap = new Map();
      const removedIds = [];
      let failedCount = 0;

      results.forEach((result, index) => {
        const server = trackedServers[index];
        if (result.status === "fulfilled") {
          nextStatusMap.set(server.server_id, result.value.status);
          return;
        }

        const errorMessage = result.reason?.message || "";
        if (/not mapped to a host|key .* not found|no such file|not found/i.test(errorMessage)) {
          removedIds.push(server.server_id);
          return;
        }

        failedCount += 1;
      });

      if (nextStatusMap.size > 0 || removedIds.length > 0) {
        setServers((previous) => previous
          .filter((server) => !removedIds.includes(server.server_id))
          .map((server) => {
            const nextStatus = nextStatusMap.get(server.server_id);
            return nextStatus ? { ...server, status: nextStatus } : server;
          }));
      }

      if (selectedServer && !removedIds.includes(selectedServer.server_id)) {
        await loadBackups(selectedServer);
      }

      if (!silent) {
        if (removedIds.length > 0) {
          pushNotice("success", `Removed ${removedIds.length} stale server${removedIds.length === 1 ? "" : "s"} from the dashboard.`);
        } else if (failedCount > 0) {
          pushNotice("error", "Some servers could not be refreshed right now.");
        } else {
          pushNotice("success", "Server list refreshed.");
        }
      }
    } finally {
      setServersRefreshing(false);
    }
  };

  const runServerAction = async (server, action) => {
    setBusy(server.server_id, true);
    try {
      if (action === "start") {
        await apiFetch("/server/start", {
          method: "POST",
          body: { server_id: server.server_id, RAM: `${server.ram}G`, cpu_cores: server.cores },
          token,
        });
        updateServer(server.server_id, () => ({ status: "online" }));
      }
      if (action === "stop") {
        await apiFetch("/server/stop", { method: "POST", body: { server_id: server.server_id }, token });
        updateServer(server.server_id, () => ({ status: "offline" }));
      }
      if (action === "delete") {
        await apiFetch("/server/delete", { method: "POST", body: { server_id: server.server_id }, token });
        setServers((previous) => previous.filter((entry) => entry.server_id !== server.server_id));
      }
      pushNotice("success", `${action[0].toUpperCase()}${action.slice(1)} request completed for ${server.name}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const refreshServerInfo = async () => {
    if (!selectedServer) return;
    await refreshOwnedServers({ silent: true });
    await loadBackups(selectedServer);
    await loadPowerUsage(selectedServer, { silent: true });
    setConsoleRefreshTick((value) => value + 1);
    pushNotice("success", "Server info refreshed.");
  };

  useEffect(() => {
    if (active !== "servers" || !token || !currentUser?.ownerKey) return undefined;

    const timer = window.setInterval(() => {
      refreshOwnedServers({ silent: true });
    }, 30000);

    return () => window.clearInterval(timer);
  }, [active, token, currentUser?.ownerKey, servers]);

  useEffect(() => {
    if (
      active !== "servers"
      || !selectedServer?.server_id
      || !token
      || !consoleWindowOpen
    ) return undefined;

    let socket = null;
    let reconnectTimer = null;
    let closed = false;
    const serverID = selectedServer.server_id;

    const connect = () => {
      if (closed) return;

      setConsoleStateMap((previous) => ({ ...previous, [serverID]: "connecting" }));
      socket = new WebSocket(buildConsoleWsUrl(serverID, token));

      socket.onopen = () => {
        setConsoleStateMap((previous) => ({ ...previous, [serverID]: "live" }));
      };

      socket.onmessage = (event) => {
        const payload = safeJsonParse(event.data, null);
        if (!payload) return;
        if (payload.error) {
          setConsoleStateMap((previous) => ({ ...previous, [serverID]: "error" }));
          return;
        }

        setConsoleTextMap((previous) => ({ ...previous, [serverID]: payload.console || "" }));
        setPlayerCountMap((previous) => ({ ...previous, [serverID]: Number(payload.online_players) || 0 }));
      };

      socket.onerror = () => {
        setConsoleStateMap((previous) => ({ ...previous, [serverID]: "error" }));
      };

      socket.onclose = () => {
        if (closed) return;
        setConsoleStateMap((previous) => ({ ...previous, [serverID]: "reconnecting" }));
        reconnectTimer = window.setTimeout(connect, 1500);
      };
    };

    connect();

    return () => {
      closed = true;
      if (reconnectTimer) window.clearTimeout(reconnectTimer);
      if (socket) socket.close();
    };
  }, [active, selectedServer?.server_id, token, consoleWindowOpen, selectedServer?.status, consoleRefreshTick]);

  useEffect(() => {
    if (!consoleWindowOpen) return;
    const element = consoleOutputRef.current;
    if (element) {
      element.scrollTop = element.scrollHeight;
    }
  }, [consoleWindowOpen, selectedServer?.server_id, consoleTextMap]);

  const createBackup = async (server) => {
    setBusy(server.server_id, true);
    try {
      await apiFetch("/server/backup/create", {
        method: "POST",
        body: { server_id: server.server_id, bundle: server.bundleName },
        token,
      });
      updateServer(server.server_id, (current) => ({ backupCount: (current.backupCount || 0) + 1 }));
      await loadBackups(server);
      pushNotice("success", `Backup created for ${server.name}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const sendConsoleCommand = async () => {
    if (!selectedServer || !consoleCommand.trim()) return;
    setConsoleSending(true);
    try {
      await apiFetch("/server/console/command", {
        method: "POST",
        body: { server_id: selectedServer.server_id, command: consoleCommand.trim() },
        token,
      });
      setConsoleCommand("");
      pushNotice("success", "Command sent.");
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setConsoleSending(false);
    }
  };

  const downloadBackup = async (server, backupID) => {
    if (!backupID) {
      pushNotice("error", "Select a backup first.");
      return;
    }

    setBusy(server.server_id, true);
    try {
      const blob = await apiFetchBlob("/server/backup/get", {
        method: "POST",
        body: { server_id: server.server_id, backup_id: backupID },
        token,
      });
      const suffix = String(backupID).slice(0, 8) || "selected";
      downloadBlob(blob, `${server.name || server.server_id}-backup-${suffix}.tar.gz`);
      pushNotice("success", `Downloaded backup ${backupID}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const loadBackups = async (server) => {
    setBackupLoading(server.server_id, true);
    try {
      const data = await apiFetch("/server/backup/list", {
        method: "POST",
        body: { server_id: server.server_id },
        token,
      });
      const nextBackups = Array.isArray(data.backups) ? data.backups : [];
      setBackupListMap((previous) => ({ ...previous, [server.server_id]: nextBackups }));
      setSelectedBackupMap((previous) => {
        const existing = previous[server.server_id];
        const hasExisting = nextBackups.some((entry) => entry.backup_id === existing);
        const nextSelected = hasExisting ? existing : nextBackups[0]?.backup_id || "";
        return { ...previous, [server.server_id]: nextSelected };
      });
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBackupLoading(server.server_id, false);
    }
  };

  const deleteBackup = async (server, backupID) => {
    setBackupLoading(server.server_id, true);
    try {
      await apiFetch("/server/backup/delete", {
        method: "POST",
        body: { server_id: server.server_id, backup_id: backupID },
        token,
      });
      await loadBackups(server);
      updateServer(server.server_id, (current) => ({ backupCount: Math.max((current.backupCount || 0) - 1, 0) }));
      pushNotice("success", `Backup ${backupID} deleted.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBackupLoading(server.server_id, false);
    }
  };

  const uploadBackup = async (server, file) => {
    if (!file) {
      pushNotice("error", "Please select a backup file");
      return;
    }
    setBusy(server.server_id, true);
    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("server_id", server.server_id);
      formData.append("backup_name", file.name);

      const response = await fetch(`${API_BASE}/server/backup/upload`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.text();
        if (response.status === 413) {
          throw new Error("Upload file is too large for the web gateway. Try a smaller archive or increase the upload limit.");
        }
        throw new Error(extractErrorMessage(errorData, `Upload failed with status ${response.status}`));
      }

      pushNotice("success", `Backup uploaded for ${server.name}`);
      updateServer(server.server_id, (current) => ({ backupCount: (current.backupCount || 0) + 1 }));
      await loadBackups(server);
    } catch (error) {
      pushNotice("error", `Backup upload failed: ${error.message}`);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const uploadWorld = async (server, file) => {
    if (!file) {
      pushNotice("error", "Please select a world folder archive");
      return;
    }
    setBusy(server.server_id, true);
    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("server_id", server.server_id);

      const response = await fetch(`${API_BASE}/server/world/upload`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.text();
        if (response.status === 413) {
          throw new Error("Upload file is too large for the web gateway. Try a smaller archive or increase the upload limit.");
        }
        throw new Error(extractErrorMessage(errorData, `Upload failed with status ${response.status}`));
      }

      pushNotice("success", `World uploaded for ${server.name}. Server will use the new world data.`);
    } catch (error) {
      pushNotice("error", `World upload failed: ${error.message}`);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const copyJoinAddress = async (server) => {
    const address = `localhost:${server.port || FIRST_SERVER_PORT}`;
    try {
      await navigator.clipboard.writeText(address);
      pushNotice("success", `Join address copied: ${address}`);
    } catch {
      pushNotice("error", `Could not copy. Join address: ${address}`);
    }
  };

  const loadPowerUsage = async (server, { silent = false } = {}) => {
    if (!server) return;

    setPowerUsageLoading((previous) => ({ ...previous, [server.server_id]: true }));
    try {
      const usage = await apiFetch(`/server/power/usage?server_id=${encodeURIComponent(server.server_id)}`, { token });
      setPowerUsageMap((previous) => ({ ...previous, [server.server_id]: usage }));
    } catch (error) {
      if (!silent) pushNotice("error", `Could not load live resource usage: ${error.message}`);
    } finally {
      setPowerUsageLoading((previous) => ({ ...previous, [server.server_id]: false }));
    }
  };

  useEffect(() => {
    if (!selectedServer || !pluginCapableServerTypes.has(selectedServer.software)) {
      setPluginResults([]);
      setPluginPage(0);
      setPluginHasMore(false);
      return;
    }

    const timer = setTimeout(async () => {
      setPluginLoading(true);
      try {
        const result = await searchModrinthPlugins(pluginSearch, pluginPage);
        setPluginResults(result.items);
        setPluginHasMore(result.hasMore);
      } catch (error) {
        pushNotice("error", `Plugin search failed: ${error.message}`);
      } finally {
        setPluginLoading(false);
      }
    }, 350);

    return () => clearTimeout(timer);
  }, [pluginSearch, pluginPage, selectedServer?.server_id, selectedServer?.software]);

  useEffect(() => {
    setPluginPage(0);
  }, [pluginSearch, selectedServer?.server_id]);

  const installPlugin = async (server, plugin) => {
    if (!server) return;

    setPluginInstalling((previous) => ({ ...previous, [plugin.id]: true }));
    try {
      const file = await fetchLatestPluginFile(plugin.id, server.version);
      const existingPlugins = installedPluginsMap[server.server_id] || [];
      if (existingPlugins.some((name) => String(name || "").toLowerCase() === String(file.name || "").toLowerCase())) {
        throw new Error(`${plugin.name} is already installed.`);
      }

      const formData = new FormData();
      formData.append("server_id", server.server_id);
      formData.append("file", file);

      const response = await fetch(`${API_BASE}/server/plugin/install`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(extractErrorMessage(errorData, `Plugin install failed with status ${response.status}`));
      }

      pushNotice("success", `${plugin.name} installed for ${server.name}. Restart the server if it is already running.`);
      await loadInstalledPlugins(server);
    } catch (error) {
      pushNotice("error", `Plugin install failed: ${error.message}`);
    } finally {
      setPluginInstalling((previous) => ({ ...previous, [plugin.id]: false }));
    }
  };

    const loadInstalledPlugins = async (server) => {
    if (!server) return;
    setInstalledPluginsLoading((previous) => ({ ...previous, [server.server_id]: true }));
    try {
      const response = await apiFetch(`/server/plugin/list?server_id=${encodeURIComponent(server.server_id)}`, { token });
      setInstalledPluginsMap((previous) => ({ ...previous, [server.server_id]: Array.isArray(response.plugins) ? response.plugins : [] }));
    } catch (error) {
      pushNotice("error", `Could not load installed plugins: ${error.message}`);
    } finally {
      setInstalledPluginsLoading((previous) => ({ ...previous, [server.server_id]: false }));
    }
    };

    const deleteInstalledPlugin = async (server, fileName) => {
    if (!server || !fileName) return;
    const key = `${server.server_id}:${fileName}`;
    setPluginDeleting((previous) => ({ ...previous, [key]: true }));
    try {
      await apiFetch("/server/plugin/delete", {
        method: "POST",
        body: { server_id: server.server_id, file_name: fileName },
        token,
      });
      pushNotice("success", `${fileName} removed from ${server.name}.`);
      await loadInstalledPlugins(server);
    } catch (error) {
      pushNotice("error", `Plugin delete failed: ${error.message}`);
    } finally {
      setPluginDeleting((previous) => ({ ...previous, [key]: false }));
    }
    };

    useEffect(() => {
    if (!selectedServer || !pluginCapableServerTypes.has(selectedServer.software) || !token) return;
    loadInstalledPlugins(selectedServer);
    }, [selectedServer?.server_id, selectedServer?.software, token]);

  useEffect(() => {
    if (!selectedServer || !token) return;
    loadBackups(selectedServer);
  }, [selectedServer?.server_id, token]);

  useEffect(() => {
    if (active !== "servers" || !selectedServer || !token) return undefined;

    loadPowerUsage(selectedServer, { silent: true });
    const timer = window.setInterval(() => {
      loadPowerUsage(selectedServer, { silent: true });
    }, 7000);

    return () => window.clearInterval(timer);
  }, [active, selectedServer?.server_id, token]);

  const menu = [
    { key: "servers", label: "My Servers", icon: Server },
    { key: "support", label: "Support", icon: MessageSquare },
    ...(SHOW_PLUGIN_CATALOG ? [{ key: "catalog", label: "Plugins & Mods", icon: Boxes }] : []),
    { key: "account", label: "Account", icon: Settings },
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-10 lg:px-16">
      <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">Server Dashboard</div>
          <h1 className="display-font mt-2 text-3xl font-bold text-white">Welcome back, {currentUser.displayName}</h1>
          <p className="mt-2 text-slate-300">Manage your Minecraft servers — start, stop, back up, and monitor from one place.</p>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <div className={cn("rounded-2xl px-4 py-3 text-sm", apiHealthy ? "bg-emerald-400/10 text-emerald-200" : "bg-amber-400/10 text-amber-100")}>
            <Wifi className="mr-2 inline h-4 w-4" />
            {apiHealthy ? "Gateway reachable" : "Gateway not reachable"}
          </div>
          <button onClick={() => onBuyServer(plans[1] || plans[0])} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-4 py-3 font-semibold text-slate-950 transition-all duration-300 hover:shadow-xl hover:shadow-cyan-400/30", popClass())}>
            Buy Server
          </button>
          <button onClick={logout} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white transition-all duration-300 hover:bg-white/10 hover:border-white/20 hover:shadow-lg hover:shadow-white/10", popClass())}>
            Log out
          </button>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <Sidebar items={menu} active={active} setActive={setActive} title="easy2host Panel" subtitle="User API area" />
        <div className="space-y-6">

          {active === "servers" && (
            <>
              <GlassCard className="p-5">
                <div className="mb-4 flex items-center justify-between gap-3">
                  <h2 className="display-font text-xl font-semibold text-white">Buy Another Server</h2>
                  <span className="text-sm text-slate-400">Choose a plan</span>
                </div>
                <div className="grid gap-3 sm:grid-cols-3">
                  {plans.map((plan) => (
                    <button
                      key={plan.name}
                      type="button"
                      onClick={() => onBuyServer(plan)}
                      className={cn(
                        "rounded-2xl border p-4 text-left transition-all duration-300 hover:shadow-xl hover:shadow-cyan-400/20",
                        popClass(),
                        plan.featured ? "border-cyan-300/30 bg-cyan-400/10 hover:border-cyan-300/50" : "border-white/10 bg-slate-950/40 hover:border-white/20 hover:bg-slate-950/60",
                      )}
                    >
                      <div className="font-semibold text-white">{plan.name}</div>
                      <div className="mt-1 text-sm text-slate-300">{plan.ram} GB RAM · {plan.cores} cores</div>
                      <div className="mt-2 text-sm text-cyan-300">{money(plan.price)} / month</div>
                    </button>
                  ))}
                </div>
              </GlassCard>

              <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                <GlassCard className={cn("p-5", popClass())}>
                  <div className="text-sm text-slate-400">Owned servers</div>
                  <div className="display-font mt-2 text-3xl font-bold text-white">{ownedServers.length}</div>
                </GlassCard>
                <GlassCard className={cn("p-5", popClass())}>
                  <div className="text-sm text-slate-400">Running</div>
                  <div className="display-font mt-2 text-3xl font-bold text-white">{ownedServers.filter((server) => server.status === "online").length}</div>
                </GlassCard>
                <GlassCard className={cn("p-5", popClass())}>
                  <div className="text-sm text-slate-400">Known backups</div>
                  <div className="display-font mt-2 text-3xl font-bold text-white">{ownedServers.reduce((sum, server) => sum + (server.backupCount || 0), 0)}</div>
                </GlassCard>
              </div>

              <div className="grid gap-6 2xl:grid-cols-[0.95fr_1.05fr]">
                <GlassCard className="p-5">
                  <div className="mb-4 flex items-center justify-between gap-4">
                    <h2 className="display-font text-xl font-semibold text-white">My Servers</h2>
                    <div className="flex items-center gap-3">
                      <button
                        type="button"
                        onClick={() => refreshOwnedServers()}
                        disabled={serversRefreshing}
                        className={cn("inline-flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-1.5 text-sm text-white disabled:opacity-60", popClass())}
                      >
                        <RefreshCw className={cn("h-4 w-4", serversRefreshing && "animate-spin")} />
                        Refresh Servers
                      </button>
                      <div className="text-sm text-slate-400">Your servers across all plans</div>
                    </div>
                  </div>
                  <div className="space-y-3">
                    {ownedServers.length === 0 && (
                      <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">
                        No tracked servers yet. Create one from the purchase flow and it will appear here immediately.
                      </div>
                    )}
                    {ownedServers.map((server) => (
                      <button
                        key={server.server_id}
                        onClick={() => setSelectedServerId(server.server_id)}
                        className={cn(
                          "w-full rounded-2xl border p-4 text-left transition-all duration-300",
                          popClass(),
                          selectedServer?.server_id === server.server_id
                            ? "border-cyan-300/30 bg-cyan-400/10 shadow-lg shadow-cyan-400/10"
                            : "border-white/10 bg-slate-950/40 hover:bg-white/5 hover:border-white/20 hover:shadow-lg hover:shadow-white/5",
                        )}
                      >
                        <div className="flex items-center justify-between gap-4">
                          <div>
                            <div className="font-semibold text-white">{server.name}</div>
                            <div className="mt-1 text-sm text-slate-400">
                              {server.software} {server.version} · {server.bundleName}
                            </div>
                          </div>
                          <div className={cn("rounded-full px-3 py-1 text-xs", server.status === "online" ? "bg-emerald-400/10 text-emerald-300" : "bg-slate-700/60 text-slate-300")}>
                            {server.status || "created"}
                          </div>
                        </div>
                      </button>
                    ))}
                  </div>
                </GlassCard>

                <GlassCard className="p-5">
                  {selectedServer ? (
                    <>
                      <div className="flex flex-wrap items-start justify-between gap-4">
                        <div>
                          <h2 className="display-font text-2xl font-bold text-white">{selectedServer.name}</h2>
                          <p className="mt-1 text-slate-400">
                            {selectedServer.software} {selectedServer.version} · {selectedServer.bundleName}
                          </p>
                          <div className="mt-3 inline-flex items-center gap-2 rounded-full border border-cyan-300/25 bg-cyan-400/10 px-3 py-1 text-sm text-cyan-200">
                            Join: localhost:{selectedServer.port || FIRST_SERVER_PORT}
                            <button type="button" onClick={() => copyJoinAddress(selectedServer)} className="rounded-lg border border-cyan-300/20 bg-cyan-300/10 p-1 text-cyan-100 hover:bg-cyan-300/20">
                              <Copy className="h-3.5 w-3.5" />
                            </button>
                          </div>
                        </div>
                        <div className={cn("rounded-full px-3 py-1 text-sm", selectedServer.status === "online" ? "bg-emerald-400/10 text-emerald-300" : "bg-slate-700/60 text-slate-300")}>
                          {selectedServer.status || "created"}
                        </div>
                        <div className="rounded-full px-3 py-1 text-sm bg-cyan-400/10 text-cyan-200">
                          Players: {playerCountMap[selectedServer.server_id] || 0}
                        </div>
                      </div>

                      <div className="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "start")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-emerald-400/15 px-4 py-3 text-emerald-300 disabled:opacity-60 transition-all duration-300 hover:bg-emerald-400/25 hover:shadow-lg hover:shadow-emerald-400/20", popClass())}><Play className="h-4 w-4" /> Start</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "stop")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-400/15 px-4 py-3 text-rose-300 disabled:opacity-60 transition-all duration-300 hover:bg-rose-400/25 hover:shadow-lg hover:shadow-rose-400/20", popClass())}><Square className="h-4 w-4" /> Stop</button>
                        <button
                          onClick={() => setConsoleWindowOpen(true)}
                          className={cn("flex items-center justify-center gap-2 rounded-2xl bg-cyan-400/15 px-4 py-3 text-cyan-200 transition-all duration-300 hover:bg-cyan-400/25 hover:shadow-lg hover:shadow-cyan-400/20", popClass())}
                        >
                          <TerminalSquare className="h-4 w-4" /> Open Console
                        </button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => createBackup(selectedServer)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-violet-400/15 px-4 py-3 text-violet-300 disabled:opacity-60 transition-all duration-300 hover:bg-violet-400/25 hover:shadow-lg hover:shadow-violet-400/20", popClass())}><HardDrive className="h-4 w-4" /> Create Backup</button>
                          <button disabled={serverBusy[selectedServer.server_id]} onClick={() => setUploadBackupModalOpen(true)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-amber-400/15 px-4 py-3 text-amber-300 disabled:opacity-60 transition-all duration-300 hover:bg-amber-400/25 hover:shadow-lg hover:shadow-amber-400/20", popClass())}><Download className="h-4 w-4" /> Upload Backup</button>
                          <button disabled={serverBusy[selectedServer.server_id]} onClick={() => setUploadWorldModalOpen(true)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-sky-400/15 px-4 py-3 text-sky-300 disabled:opacity-60 transition-all duration-300 hover:bg-sky-400/25 hover:shadow-lg hover:shadow-sky-400/20", popClass())}><HardDrive className="h-4 w-4" /> Upload World</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "delete")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-500/20 px-4 py-3 text-rose-200 disabled:opacity-60 transition-all duration-300 hover:bg-rose-500/30 hover:shadow-lg hover:shadow-rose-500/20", popClass())}><Trash2 className="h-4 w-4" /> Delete</button>
                      </div>

                      <div className="mt-6">
                        <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex items-center justify-between gap-2 text-white">
                            <div className="flex items-center gap-2"><Database className="h-4 w-4 text-cyan-300" /> Server Details</div>
                            <button
                              type="button"
                              onClick={refreshServerInfo}
                              disabled={backupBusy[selectedServer.server_id]}
                              className={cn("rounded-xl border border-white/10 bg-white/5 px-3 py-1.5 text-sm text-white disabled:opacity-60", popClass())}
                            >
                              Refresh Info
                            </button>
                          </div>
                          <div className="space-y-2 text-sm text-slate-300">
                            <div className="flex justify-between"><span className="text-slate-400">Plan</span><span className="text-white">{selectedServer.bundleName}</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">RAM</span><span className="text-white">{selectedServer.ram} GB</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">Storage</span><span className="text-white">{selectedServer.storage} GB SSD</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">CPU cores</span><span className="text-white">{selectedServer.cores}</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">Software</span><span className="text-white">{selectedServer.software} {selectedServer.version}</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">Port</span><span className="text-white">{selectedServer.port || "—"}</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">Players online</span><span className="text-cyan-200 font-semibold">{playerCountMap[selectedServer.server_id] || 0}</span></div>
                            <div className="flex justify-between"><span className="text-slate-400">Backups</span><span className="text-white">{selectedServer.backupCount || 0}</span></div>
                          </div>

                          <div className="mt-4 rounded-xl border border-white/10 bg-black/20 p-3">
                            <div className="mb-2 flex items-center justify-between gap-2">
                              <div className="text-sm font-semibold text-white">Live Power Usage</div>
                              {powerUsageLoading[selectedServer.server_id] && <div className="text-xs text-slate-400">Updating...</div>}
                            </div>

                            {(() => {
                              const usage = powerUsageMap[selectedServer.server_id] || {};
                              const cpuPercent = clampPercent(Number(usage.cpu_percent) || 0);
                              const ramPercent = clampPercent(Number(usage.ram_percent) || 0);
                              const ramUsedMB = Number(usage.ram_used_mb) || 0;
                              const ramLimitMB = Number(usage.ram_limit_mb) || 0;
                              const online = Boolean(usage.online);

                              return (
                                <div className="space-y-3">
                                  <div className="text-xs text-slate-400">
                                    {online ? "Live usage from the running container" : "Server is offline. Usage will appear after start."}
                                  </div>

                                  <div>
                                    <div className="mb-1 flex justify-between text-xs text-slate-300">
                                      <span>CPU</span>
                                      <span>{cpuPercent.toFixed(1)}%</span>
                                    </div>
                                    <div className="h-2 overflow-hidden rounded-full bg-slate-800">
                                      <div
                                        className="h-full rounded-full bg-[repeating-linear-gradient(45deg,rgba(56,189,248,0.95),rgba(56,189,248,0.95)_8px,rgba(14,116,144,0.95)_8px,rgba(14,116,144,0.95)_16px)] transition-all duration-700"
                                        style={{ width: `${cpuPercent}%` }}
                                      />
                                    </div>
                                  </div>

                                  <div>
                                    <div className="mb-1 flex justify-between text-xs text-slate-300">
                                      <span>RAM</span>
                                      <span>
                                        {ramUsedMB.toFixed(0)} MB / {ramLimitMB > 0 ? `${ramLimitMB.toFixed(0)} MB` : "n/a"} ({ramPercent.toFixed(1)}%)
                                      </span>
                                    </div>
                                    <div className="h-2 overflow-hidden rounded-full bg-slate-800">
                                      <div
                                        className="h-full rounded-full bg-[repeating-linear-gradient(45deg,rgba(74,222,128,0.95),rgba(74,222,128,0.95)_8px,rgba(22,163,74,0.95)_8px,rgba(22,163,74,0.95)_16px)] transition-all duration-700"
                                        style={{ width: `${ramPercent}%` }}
                                      />
                                    </div>
                                  </div>
                                </div>
                              );
                            })()}
                          </div>
                        </div>

                        <div className="mt-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
                            <div className="flex items-center gap-2 text-white"><HardDrive className="h-4 w-4 text-cyan-300" /> Backup List</div>
                            <button
                              type="button"
                              onClick={() => loadBackups(selectedServer)}
                              disabled={backupBusy[selectedServer.server_id]}
                              className={cn("rounded-xl border border-white/10 bg-white/5 px-3 py-1.5 text-sm text-white disabled:opacity-60", popClass())}
                            >
                              Refresh
                            </button>
                          </div>

                          {(backupListMap[selectedServer.server_id] || []).length === 0 && (
                            <div className="rounded-xl border border-dashed border-white/10 bg-black/20 px-3 py-3 text-sm text-slate-400">
                              No backups recorded yet.
                            </div>
                          )}

                          <div className="space-y-2">
                            {(backupListMap[selectedServer.server_id] || []).map((entry) => (
                              <div key={entry.backup_id} className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-white/10 bg-black/20 px-3 py-2">
                                <button
                                  type="button"
                                  onClick={() => setSelectedBackupMap((previous) => ({ ...previous, [selectedServer.server_id]: entry.backup_id }))}
                                  className={cn(
                                    "min-w-0 rounded-lg border px-2.5 py-1 text-left text-sm",
                                    selectedBackupMap[selectedServer.server_id] === entry.backup_id
                                      ? "border-cyan-300/35 bg-cyan-400/10 text-cyan-200"
                                      : "border-white/10 bg-white/5 text-slate-300",
                                  )}
                                >
                                  {entry.backup_id}
                                </button>
                                <div className="flex items-center gap-2">
                                  <button
                                    type="button"
                                    disabled={serverBusy[selectedServer.server_id]}
                                    onClick={() => downloadBackup(selectedServer, entry.backup_id)}
                                    className={cn("rounded-lg border border-cyan-400/25 bg-cyan-400/10 px-2.5 py-1 text-xs font-semibold text-cyan-200 disabled:opacity-60", popClass())}
                                  >
                                    Download
                                  </button>
                                  <button
                                    type="button"
                                    disabled={backupBusy[selectedServer.server_id]}
                                    onClick={() => deleteBackup(selectedServer, entry.backup_id)}
                                    className={cn("rounded-lg border border-rose-400/25 bg-rose-400/10 px-2.5 py-1 text-xs font-semibold text-rose-200 disabled:opacity-60", popClass())}
                                  >
                                    Delete
                                  </button>
                                </div>
                              </div>
                            ))}
                          </div>
                        </div>

                        {pluginCapableServerTypes.has(selectedServer.software) && (
                          <div className="mt-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                            <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                              <div>
                                <div className="flex items-center gap-2 text-white"><Boxes className="h-4 w-4 text-cyan-300" /> Plugins</div>
                                <p className="mt-1 text-sm text-slate-400">Browse Modrinth plugins compatible with {selectedServer.software} and install them with one click.</p>
                              </div>
                              <div className="relative w-full max-w-sm">
                                <Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" />
                                <input
                                  value={pluginSearch}
                                  onChange={(event) => setPluginSearch(event.target.value)}
                                  className="w-full rounded-2xl border border-white/10 bg-white/5 py-3 pl-11 pr-4 text-white outline-none"
                                  placeholder="Search Modrinth plugins"
                                />
                              </div>
                            </div>

                            <div className="mb-4 flex items-center justify-end gap-2">
                              <button
                                type="button"
                                onClick={() => setPluginPage((previous) => Math.max(previous - 1, 0))}
                                disabled={pluginLoading || pluginPage === 0}
                                className={cn("rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-semibold text-slate-200 disabled:opacity-50", popClass())}
                                aria-label="Previous plugin page"
                              >
                                <ChevronLeft className="h-4 w-4" />
                              </button>
                              <div className="text-xs text-slate-400">Page {pluginPage + 1}</div>
                              <button
                                type="button"
                                onClick={() => setPluginPage((previous) => previous + 1)}
                                disabled={pluginLoading || !pluginHasMore}
                                className={cn("rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-semibold text-slate-200 disabled:opacity-50", popClass())}
                                aria-label="Next plugin page"
                              >
                                <ChevronRight className="h-4 w-4" />
                              </button>
                            </div>

                            {pluginLoading ? (
                              <div className="rounded-xl border border-dashed border-white/10 bg-black/20 px-3 py-3 text-sm text-slate-400">
                                Loading plugins...
                              </div>
                            ) : pluginResults.length === 0 ? (
                              <div className="rounded-xl border border-dashed border-white/10 bg-black/20 px-3 py-3 text-sm text-slate-400">
                                No plugins found.
                              </div>
                            ) : (
                              <div className="grid gap-3 xl:grid-cols-2">
                                {pluginResults.map((plugin) => {
                                  const installedFiles = installedPluginsMap[selectedServer.server_id] || [];
                                  const alreadyInstalled = isPluginAlreadyInstalled(plugin, installedFiles);
                                  return (
                                    <div key={plugin.id} className="rounded-xl border border-white/10 bg-black/20 p-4">
                                      <div className="flex items-start justify-between gap-3">
                                        <div className="flex min-w-0 items-start gap-3">
                                          {plugin.iconUrl ? (
                                            <img
                                              src={plugin.iconUrl}
                                              alt={`${plugin.name} icon`}
                                              className="h-10 w-10 shrink-0 rounded-lg border border-white/10 object-cover bg-slate-900"
                                              loading="lazy"
                                            />
                                          ) : (
                                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-white/10 bg-slate-900 text-xs font-semibold text-cyan-200">
                                              {String(plugin.name || "P").charAt(0).toUpperCase()}
                                            </div>
                                          )}
                                          <div className="min-w-0">
                                            <div className="truncate font-semibold text-white">{plugin.name}</div>
                                            <div className="mt-1 text-xs text-slate-400">by {plugin.author} · {typeof plugin.downloads === "number" ? plugin.downloads.toLocaleString() : plugin.downloads || 0} downloads</div>
                                          </div>
                                        </div>
                                        <button
                                          type="button"
                                          onClick={() => installPlugin(selectedServer, plugin)}
                                          disabled={pluginInstalling[plugin.id] || alreadyInstalled}
                                          className={cn("rounded-lg border border-cyan-400/25 bg-cyan-400/10 px-3 py-1.5 text-xs font-semibold text-cyan-200 disabled:opacity-60", popClass())}
                                        >
                                          {pluginInstalling[plugin.id] ? "Installing..." : alreadyInstalled ? "Installed" : "Install"}
                                        </button>
                                      </div>
                                      <p className="mt-3 text-sm leading-6 text-slate-300">{plugin.description}</p>
                                    </div>
                                  );
                                })}
                              </div>
                            )}

                            <div className="mt-4 rounded-xl border border-white/10 bg-black/20 p-3">
                              <div className="mb-2 flex items-center justify-between gap-2">
                                <div className="text-sm font-semibold text-white">Installed Plugins</div>
                                <button
                                  type="button"
                                  onClick={() => loadInstalledPlugins(selectedServer)}
                                  disabled={installedPluginsLoading[selectedServer.server_id]}
                                  className={cn("rounded-lg border border-white/10 bg-white/5 px-2.5 py-1 text-xs text-slate-200 disabled:opacity-60", popClass())}
                                >
                                  Refresh
                                </button>
                              </div>

                              {installedPluginsLoading[selectedServer.server_id] ? (
                                <div className="text-xs text-slate-400">Loading installed plugins...</div>
                              ) : (installedPluginsMap[selectedServer.server_id] || []).length === 0 ? (
                                <div className="text-xs text-slate-400">No installed plugins found.</div>
                              ) : (
                                <div className="space-y-2">
                                  {(installedPluginsMap[selectedServer.server_id] || []).map((fileName) => {
                                    const deleteKey = `${selectedServer.server_id}:${fileName}`;
                                    return (
                                      <div key={fileName} className="flex items-center justify-between gap-2 rounded-lg border border-white/10 bg-slate-900/40 px-2.5 py-2">
                                        <div className="min-w-0 truncate text-xs text-slate-200">{fileName}</div>
                                        <button
                                          type="button"
                                          onClick={() => deleteInstalledPlugin(selectedServer, fileName)}
                                          disabled={pluginDeleting[deleteKey]}
                                          className={cn("rounded-lg border border-rose-400/25 bg-rose-400/10 px-2 py-1 text-xs font-semibold text-rose-200 disabled:opacity-60", popClass())}
                                        >
                                          {pluginDeleting[deleteKey] ? "Deleting..." : "Delete"}
                                        </button>
                                      </div>
                                    );
                                  })}
                                </div>
                              )}
                            </div>
                          </div>
                        )}
                      </div>
                    </>
                  ) : (
                    <div className="text-slate-300">No server selected.</div>
                  )}
                </GlassCard>
              </div>

              <Modal
                open={consoleWindowOpen && Boolean(selectedServer)}
                onClose={() => setConsoleWindowOpen(false)}
                title={selectedServer ? `${selectedServer.name} Console` : "Server Console"}
                subtitle="Live output and command input"
              >
                {selectedServer ? (
                <>
                    <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
                      <div className="rounded-full px-3 py-1 text-xs bg-cyan-400/10 text-cyan-200">Players online: {playerCountMap[selectedServer.server_id] || 0}</div>
                      <div className={cn(
                        "rounded-full px-3 py-1 text-xs",
                        consoleStateMap[selectedServer.server_id] === "live"
                          ? "bg-emerald-400/10 text-emerald-300"
                          : consoleStateMap[selectedServer.server_id] === "error"
                            ? "bg-rose-400/10 text-rose-300"
                            : "bg-amber-400/10 text-amber-200",
                      )}>
                        {consoleStateMap[selectedServer.server_id] || "connecting"}
                      </div>
                    </div>
                    <pre ref={consoleOutputRef} className="max-h-[50vh] overflow-auto whitespace-pre-wrap rounded-xl border border-white/10 bg-black/40 p-3 text-xs text-slate-200">{consoleTextMap[selectedServer.server_id] || "Waiting for console output..."}</pre>
                    <div className="mt-4 flex gap-2">
                      <input
                        value={consoleCommand}
                        onChange={(event) => setConsoleCommand(event.target.value)}
                        onKeyDown={(event) => {
                          if (event.key === "Enter" && !consoleSending) {
                            event.preventDefault();
                            sendConsoleCommand();
                          }
                        }}
                        placeholder="Type command (example: say Hello players)"
                        className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
                      />
                      <button
                        type="button"
                        disabled={consoleSending || !consoleCommand.trim()}
                        onClick={sendConsoleCommand}
                        className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}
                      >
                        Send
                      </button>
                    </div>
                </>
                ) : null}
              </Modal>
            </>
          )}

          {active === "catalog" && SHOW_PLUGIN_CATALOG && <PluginCatalogPanel search={search} setSearch={setSearch} />}
              <Modal
                open={uploadBackupModalOpen && Boolean(selectedServer)}
                onClose={() => setUploadBackupModalOpen(false)}
                title={selectedServer ? `${selectedServer.name} - Upload Backup` : "Upload Backup"}
                subtitle="Select a .tar.gz or .tar backup file from Easy2Host"
              >
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="block text-sm font-medium text-slate-300">Backup File</label>
                    <input
                      type="file"
                      ref={uploadBackupFileRef}
                      accept=".tar,.tar.gz,.tgz"
                      className="block w-full text-sm text-slate-400 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-xs file:font-semibold file:bg-cyan-400/20 file:text-cyan-300 hover:file:bg-cyan-400/30"
                    />
                    <p className="text-xs text-slate-400">Supported formats: .tar, .tar.gz, .tgz</p>
                  </div>
                  <div className="flex gap-2 pt-2">
                    <button
                      onClick={async () => {
                        if (!selectedServer) {
                          pushNotice({ type: "error", message: "No server selected" });
                          return;
                        }
                        const file = uploadBackupFileRef.current?.files?.[0];
                        if (!file) {
                          pushNotice({ type: "error", message: "Please select a backup file" });
                          return;
                        }
                        if (!file.name.endsWith('.tar') && !file.name.endsWith('.tar.gz') && !file.name.endsWith('.tgz')) {
                          pushNotice({ type: "error", message: "File must be .tar, .tar.gz, or .tgz" });
                          return;
                        }
                        await uploadBackup(selectedServer, file);
                        setUploadBackupModalOpen(false);
                        uploadBackupFileRef.current.value = '';
                      }}
                      className="flex-1 rounded-lg bg-gradient-to-r from-cyan-400 to-blue-500 px-4 py-2 font-semibold text-slate-950 hover:opacity-90 transition-opacity"
                    >
                      Upload Backup
                    </button>
                    <button
                      onClick={() => {
                        setUploadBackupModalOpen(false);
                        uploadBackupFileRef.current.value = '';
                      }}
                      className="flex-1 rounded-lg border border-slate-500 px-4 py-2 font-semibold text-slate-300 hover:bg-slate-500/10 transition-colors"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </Modal>

              <Modal
                open={uploadWorldModalOpen && Boolean(selectedServer)}
                onClose={() => setUploadWorldModalOpen(false)}
                title={selectedServer ? `${selectedServer.name} - Upload World` : "Upload World"}
                subtitle="Select a world folder archive (.tar.gz or .tar)"
              >
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="block text-sm font-medium text-slate-300">World Folder Archive</label>
                    <input
                      type="file"
                      ref={uploadWorldFileRef}
                      accept=".tar,.tar.gz,.tgz"
                      className="block w-full text-sm text-slate-400 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-xs file:font-semibold file:bg-sky-400/20 file:text-sky-300 hover:file:bg-sky-400/30"
                    />
                    <p className="text-xs text-slate-400">Supported formats: .tar, .tar.gz, .tgz - folder will replace existing world</p>
                  </div>
                  <div className="flex gap-2 pt-2">
                    <button
                      onClick={() => {
                        if (!selectedServer) {
                          pushNotice("error", "No server selected");
                          return;
                        }
                        const file = uploadWorldFileRef.current?.files?.[0];
                        if (!file) {
                          pushNotice("error", "Please select a world file");
                          return;
                        }
                        if (!file.name.endsWith('.tar') && !file.name.endsWith('.tar.gz') && !file.name.endsWith('.tgz')) {
                          pushNotice("error", "File must be .tar, .tar.gz, or .tgz");
                          return;
                        }
                        setUploadWorldModalOpen(false);
                        uploadWorldFileRef.current.value = '';
                        pushNotice("info", `Uploading world for ${selectedServer.name}...`);
                        uploadWorld(selectedServer, file);
                      }}
                      className="flex-1 rounded-lg bg-gradient-to-r from-sky-400 to-blue-500 px-4 py-2 font-semibold text-slate-950 hover:opacity-90 transition-opacity"
                    >
                      Upload World
                    </button>
                    <button
                      onClick={() => {
                        setUploadWorldModalOpen(false);
                        uploadWorldFileRef.current.value = '';
                      }}
                      className="flex-1 rounded-lg border border-slate-500 px-4 py-2 font-semibold text-slate-300 hover:bg-slate-500/10 transition-colors"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              </Modal>

          {active === "support" && (
            <UserSupportPanel
              currentUser={currentUser}
              tickets={tickets}
              setTickets={setTickets}
              pushNotice={pushNotice}
            />
          )}

          {active === "account" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Account</div>
              <h2 className="display-font mt-2 text-2xl font-bold text-white">Profile & Session</h2>
              <div className="mt-8 grid gap-4 sm:grid-cols-2 md:grid-cols-4">
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}><div className="text-sm text-slate-400">Display name</div><div className="mt-2 font-semibold text-white truncate">{currentUser.displayName}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}><div className="text-sm text-slate-400">Username</div><div className="mt-2 font-semibold text-white truncate">{currentUser.username}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}><div className="text-sm text-slate-400">Email</div><div className="mt-2 font-semibold text-white truncate">{currentUser.email || "Not provided"}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}><div className="text-sm text-slate-400">Role</div><div className="mt-2 inline-flex items-center gap-2 font-semibold text-white"><Lock className="h-4 w-4 text-cyan-300" /> {currentUser.role}</div></div>
              </div>
              <div className="mt-6 grid gap-4 sm:grid-cols-2">
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}>
                  <div className="mb-3 flex items-center gap-2 text-white"><KeyRound className="h-4 w-4 text-cyan-300" /> Session</div>
                  <div className="text-sm text-slate-300">You are signed in as <span className="text-white font-semibold">{currentUser.username}</span>.</div>
                  <button
                    type="button"
                    onClick={logout}
                    className="mt-4 w-full rounded-2xl border border-rose-400/20 bg-rose-400/10 px-4 py-2.5 text-sm font-semibold text-rose-200 transition-all duration-300 hover:bg-rose-400/20 hover:shadow-lg hover:shadow-rose-400/20"
                  >
                    Sign out
                  </button>
                </div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4 transition-all duration-300 hover:border-cyan-300/20 hover:bg-slate-950/60 hover:shadow-lg hover:shadow-cyan-400/10", popClass())}>
                  <div className="mb-3 flex items-center gap-2 text-white"><Shield className="h-4 w-4 text-cyan-300" /> Auth Token</div>
                  <div className="text-xs text-slate-400 font-mono break-all leading-relaxed">{token.slice(0, 20)}...{token.slice(-20)}</div>
                  <button
                    type="button"
                    onClick={() => {
                      navigator.clipboard.writeText(token);
                      pushNotice("success", "Token copied to clipboard");
                    }}
                    className="mt-3 w-full rounded-lg border border-cyan-400/25 bg-cyan-400/10 px-3 py-2 text-xs font-semibold text-cyan-200 transition-all duration-300 hover:bg-cyan-400/20 hover:shadow-lg hover:shadow-cyan-400/10"
                  >
                    Copy Token
                  </button>
                </div>
              </div>
            </GlassCard>
          )}
        </div>
      </div>
    </div>
  );
}

function HostMetadataTable({ metadata }) {
  const rows = Array.isArray(metadata?.metadata) ? metadata.metadata : [];

  if (rows.length === 0) {
    return <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">No host metadata loaded yet.</div>;
  }

  return (
    <div className="space-y-3">
      {rows.map((row) => (
        <div key={row.host_server_id} className="rounded-2xl border border-white/10 bg-slate-950/40 p-4 text-sm text-slate-300">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="font-semibold text-white">{row.host_server_id}</div>
            <div className="text-xs uppercase tracking-[0.18em] text-slate-500">{row.created_at}</div>
          </div>
          <div className="mt-3 grid gap-2 sm:grid-cols-3">
            <div>RAM: <span className="text-white">{row.ram ?? row.RAM ?? "n/a"}</span></div>
            <div>Cores: <span className="text-white">{row.cpu_cores ?? row.cores ?? "n/a"}</span></div>
            <div>Servers: <span className="text-white">{Object.keys(row.servers || {}).length}</span></div>
          </div>
          <pre className="mt-4 overflow-auto rounded-2xl bg-black/20 p-3 text-xs text-slate-400">{JSON.stringify(row, null, 2)}</pre>
        </div>
      ))}
    </div>
  );
}

function AdminDashboard({ currentUser, token, notices, setNotices, apiHealthy, logout, tickets, setTickets }) {
  const [active, setActive] = useState("overview");
  const [networkIp, setNetworkIp] = useState("");
  const [networkHostId, setNetworkHostId] = useState("");
  const [networkResult, setNetworkResult] = useState(null);
  const [hostMetadata, setHostMetadata] = useState(null);
  const [hostCreateForm, setHostCreateForm] = useState({ ram: "", cores: "" });
  const [hostDeleteId, setHostDeleteId] = useState("");
  const [hostAddForm, setHostAddForm] = useState({ host_server_id: "", server_id: "" });
  const [busy, setBusy] = useState(false);

  const metadataRows = Array.isArray(hostMetadata?.metadata) ? hostMetadata.metadata : [];

  const pushNotice = (type, text) => setNotices((previous) => [{ id: Date.now() + Math.random(), type, text }, ...previous].slice(0, 5));

  const callAdmin = async (label, fn) => {
    setBusy(true);
    try {
      const result = await fn();
      pushNotice("success", `${label} completed successfully.`);
      return result;
    } catch (error) {
      pushNotice("error", error.message);
      return null;
    } finally {
      setBusy(false);
    }
  };

  const adminMenu = [
    { key: "overview", label: "Overview", icon: Crown },
    { key: "tickets", label: "Support Tickets", icon: MessageSquare },
    { key: "network", label: "Network", icon: Server },
    { key: "metadata", label: "Host Metadata", icon: Database },
    { key: "mapping", label: "Add Server To Host", icon: Plus },
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-10 lg:px-16">
      <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">Admin Dashboard</div>
          <h1 className="display-font mt-2 text-3xl font-bold text-white">Infrastructure Control</h1>
          <p className="mt-2 text-slate-300">Manage your hosting infrastructure, register new hosts, and monitor your network.</p>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <div className={cn("rounded-2xl px-4 py-3 text-sm", apiHealthy ? "bg-emerald-400/10 text-emerald-200" : "bg-amber-400/10 text-amber-100")}>
            <Wifi className="mr-2 inline h-4 w-4" />
            {apiHealthy ? "Gateway reachable" : "Gateway not reachable"}
          </div>
          <button onClick={logout} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white", popClass())}>
            Log out
          </button>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <Sidebar items={adminMenu} active={active} setActive={setActive} title="Admin API Area" subtitle="Documented admin actions" />
        <div className="space-y-6">

          {active === "overview" && (
            <div className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Signed in as</div><div className="mt-2 text-lg font-bold text-white truncate">{currentUser.username}</div></GlassCard>
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Role</div><div className="display-font mt-2 text-3xl font-bold text-white">Admin</div></GlassCard>
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Session</div><div className="display-font mt-2 text-3xl font-bold text-white">{token ? "Active" : "None"}</div></GlassCard>
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Open tickets</div><div className="display-font mt-2 text-3xl font-bold text-white">{tickets.filter((ticket) => ticket.status !== "resolved").length}</div></GlassCard>
              </div>
              <GlassCard className="p-6">
                <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Admin Guide</div>
                <h2 className="display-font mt-2 text-2xl font-bold text-white">How to use the admin dashboard</h2>
                <div className="mt-6 grid gap-4 sm:grid-cols-2">
                  {[
                    ["1. Create host metadata", "Go to Host Metadata tab → enter RAM and CPU cores → click Create Metadata. Copy the generated host ID."],
                    ["2. Register the host network", "Go to Network tab → enter the host machine's IP and paste the host ID from step 1 → click Create."],
                    ["3. Users can now buy plans", "When a user buys a plan, servers are automatically assigned to an available host from your metadata registry."],
                    ["4. Handle support tickets", "Go to Support Tickets tab to read user messages, update their status (Open / In Progress / Resolved), and reply."],
                    ["5. Manage hosts", "Use Add Server To Host to manually attach a specific server ID to a host for tracking and routing purposes."],
                    ["6. Admin accounts", "Use your pre-created admin credentials to sign in and access this dashboard."],
                  ].map(([title, text]) => (
                    <div key={title} className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                      <div className="font-semibold text-white">{title}</div>
                      <p className="mt-2 text-sm leading-6 text-slate-300">{text}</p>
                    </div>
                  ))}
                </div>
              </GlassCard>
            </div>
          )}

          {active === "tickets" && <AdminTicketPanel tickets={tickets} setTickets={setTickets} pushNotice={pushNotice} />}

          {active === "network" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Admin</div>
              <h2 className="display-font mt-2 text-2xl font-bold text-white">Create Network Host Entry</h2>
              <p className="mt-3 text-slate-300">Register a host machine's IP address. You must create a host metadata entry first — the host ID from that step is required here.</p>
              <div className="mt-6 grid gap-4 sm:grid-cols-2">
                <input value={networkIp} onChange={(event) => setNetworkIp(event.target.value)} placeholder="Host IP (example: host.docker.internal or 192.168.1.20)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                <input value={networkHostId} onChange={(event) => setNetworkHostId(event.target.value)} placeholder="Host server ID (example: d2f0d4f9-3b74-4d2b-a8c7-30f9d7c8a123)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
              </div>
              {metadataRows.length > 0 && (
                <div className="mt-4">
                  <label className="mb-2 block text-xs uppercase tracking-[0.16em] text-slate-400">Pick existing host metadata ID</label>
                  <select
                    value={networkHostId}
                    onChange={(event) => setNetworkHostId(event.target.value)}
                    className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none"
                  >
                    <option value="">Select host metadata ID</option>
                    {metadataRows.map((row) => (
                      <option key={row.host_server_id} value={row.host_server_id}>{row.host_server_id}</option>
                    ))}
                  </select>
                </div>
              )}
              <div className="mt-4 flex">
                <button disabled={busy} onClick={async () => {
                  if (!networkIp.trim()) {
                    pushNotice("error", "Host IP is required.");
                    return;
                  }
                  if (!networkHostId.trim()) {
                    pushNotice("error", "Host metadata ID is required. Create metadata first and use that ID.");
                    return;
                  }

                  const payload = { ip: networkIp.trim(), host_server_id: networkHostId.trim() };
                  const result = await callAdmin("Network create", () => apiFetch("/network/create", { method: "POST", body: payload, token }));
                  if (result) setNetworkResult(result);
                }} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Create</button>
              </div>
              <div className="mt-6 rounded-2xl border border-white/10 bg-slate-950/40 p-4"><pre className="overflow-auto text-xs text-slate-300">{JSON.stringify(networkResult || { message: "No result yet." }, null, 2)}</pre></div>
            </GlassCard>
          )}

          {active === "metadata" && (
            <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
              <GlassCard className="p-6">
                <h2 className="display-font text-2xl font-bold text-white">Host Metadata</h2>
                <p className="mt-3 text-slate-300">Create, fetch, or delete host metadata through the admin gateway.</p>
                <div className="mt-6 grid gap-4 md:grid-cols-2">
                  <input value={hostCreateForm.ram} onChange={(event) => setHostCreateForm((current) => ({ ...current, ram: event.target.value }))} placeholder="RAM (example: 16)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostCreateForm.cores} onChange={(event) => setHostCreateForm((current) => ({ ...current, cores: event.target.value }))} placeholder="CPU cores (example: 8)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostDeleteId} onChange={(event) => setHostDeleteId(event.target.value)} placeholder="Host server ID to delete (example: d2f0...a123)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none md:col-span-2" />
                </div>
                <div className="mt-6 flex flex-wrap gap-3">
                  <button disabled={busy} onClick={async () => {
                    const result = await callAdmin("Metadata create", () => apiFetch("/host-metadata/create", { method: "POST", body: hostCreateForm, token }));
                    if (!result?.host_server_id) {
                      return;
                    }

                    const nextHostID = result.host_server_id;
                    setHostDeleteId(nextHostID);
                    setNetworkHostId(nextHostID);
                    setHostAddForm((current) => ({ ...current, host_server_id: nextHostID }));
                    setHostMetadata((current) => {
                      const currentRows = Array.isArray(current?.metadata) ? current.metadata : [];
                      const exists = currentRows.some((row) => row.host_server_id === nextHostID);
                      if (exists) return current;

                      return {
                        metadata: [
                          ...currentRows,
                          {
                            host_server_id: nextHostID,
                            ram: hostCreateForm.ram,
                            cpu_cores: hostCreateForm.cores,
                            servers: {},
                            created_at: new Date().toISOString(),
                          },
                        ],
                      };
                    });
                    pushNotice("success", `Host metadata created and selected: ${nextHostID}`);
                  }} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-4 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Create Metadata</button>
                  <button disabled={busy} onClick={async () => {
                    const result = await callAdmin("Metadata get", () => apiFetch("/host-metadata/get", { method: "GET", token }));
                    if (result) {
                      setHostMetadata(result);
                      const firstHostID = Array.isArray(result?.metadata) && result.metadata[0]?.host_server_id
                        ? result.metadata[0].host_server_id
                        : "";
                      if (firstHostID && !networkHostId.trim()) {
                        setNetworkHostId(firstHostID);
                        setHostAddForm((current) => ({ ...current, host_server_id: firstHostID }));
                      }
                    }
                  }} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-3 font-semibold text-white disabled:opacity-60", popClass())}>Get Metadata</button>
                  <button disabled={busy} onClick={async () => {
                    const result = await callAdmin("Metadata delete", () => apiFetch("/host-metadata/delete", { method: "POST", body: { host_server_id: hostDeleteId }, token }));
                    if (result !== null) setHostMetadata(null);
                  }} className={cn("rounded-2xl border border-rose-400/20 bg-rose-400/10 px-4 py-3 font-semibold text-rose-100 disabled:opacity-60", popClass())}>Delete Metadata</button>
                </div>
              </GlassCard>
              <GlassCard className="p-6">
                <div className="mb-4 flex items-center gap-2 text-white"><Database className="h-4 w-4 text-cyan-300" /> Metadata Response</div>
                <HostMetadataTable metadata={hostMetadata} />
              </GlassCard>
            </div>
          )}

          {active === "mapping" && (
            <GlassCard className="p-6">
              <h2 className="display-font text-2xl font-bold text-white">Add Server To Host</h2>
              <p className="mt-3 text-slate-300">Calls POST /host-metadata/add with both host and server IDs.</p>
              <div className="mt-6 grid gap-4 md:grid-cols-2">
                <input value={hostAddForm.host_server_id} onChange={(event) => setHostAddForm((current) => ({ ...current, host_server_id: event.target.value }))} placeholder="Host server ID (example: d2f0...a123)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                <input value={hostAddForm.server_id} onChange={(event) => setHostAddForm((current) => ({ ...current, server_id: event.target.value }))} placeholder="Server ID (example: b58b4f1a-98ae-4833-a866-d8cd1989d377)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
              </div>
              <div className="mt-6 flex justify-end">
                <button disabled={busy} onClick={() => callAdmin("Host mapping", () => apiFetch("/host-metadata/add", { method: "POST", body: hostAddForm, token }))} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Attach Server</button>
              </div>
            </GlassCard>
          )}
        </div>
      </div>
    </div>
  );
}

function LandingPage({ apiHealthy, currentUser, setScreen, startPurchase }) {
  const featureCards = SHOW_PLUGIN_CATALOG
    ? [
        [LayoutDashboard, "Clean dashboard", "One place for status, actions, backups, and account access."],
        [Server, "Easy server control", "Start, stop, back up, and inspect your server from one page."],
        [Boxes, "Plugins and mods", "Live catalog support is already preserved in code for later rollout."],
        [Shield, "Account-based access", "Each account only sees and manages its own infrastructure."],
        [Settings, "Version and software choice", "Pick Vanilla, Spigot, or Paper with the Minecraft version you want."],
        [CreditCard, "Simple purchase flow", "Bundle, billing, and provisioning move through one sequence."],
      ]
    : [
        [LayoutDashboard, "Clean dashboard", "One place for status, actions, backups, and settings."],
        [Server, "Full server control", "Start, stop, back up, and monitor from one clean page."],
        [HardDrive, "Backup anytime", "Create a backup and download a copy directly from your dashboard."],
        [Shield, "Your servers, your data", "Each account only sees and controls its own Minecraft servers."],
        [Settings, "Version and software choice", "Pick Vanilla, Spigot, or Paper with any supported Minecraft version."],
        [Crown, "Admin oversight", "Admin accounts give access to infrastructure management through the same sign-in."],
      ];

  const faqItems = [
    ["How fast can I create a server?", "Once your account is ready, you can choose a plan, finish checkout, and provision through the existing backend in a few steps."],
    ["Which server software can I choose?", "You can choose Vanilla, Spigot, or Paper during setup while still selecting the exact Minecraft version."],
    ["Can I add extra resources?", "Yes. During checkout you can add extra RAM and SSD upgrades before the server is created."],
    ["Will other users be able to see my server?", "No. Your dashboard only tracks the servers associated with your signed-in account."],
    ["Can I install plugins and mods?", "Plugin and mod support is coming soon. We're curating a marketplace of essential additions for your server."],
    ["How do I get support?", "Visit the Support section in your dashboard to open tickets. Our team responds to all inquiries within 24 hours."],
  ];

  return (
    <>
      <section className="relative overflow-hidden px-4 pb-16 pt-8 sm:px-6 sm:pb-20 sm:pt-10 lg:px-16 lg:pb-28 lg:pt-16">
        <div className="absolute inset-0 -z-10">
          <div className="absolute left-1/2 top-[-20%] h-[28rem] w-[28rem] -translate-x-1/2 rounded-full bg-cyan-500/15 blur-3xl" />
          <div className="absolute right-[10%] top-[15%] h-[18rem] w-[18rem] rounded-full bg-blue-600/20 blur-3xl" />
          <div className="absolute left-[5%] top-[35%] h-[16rem] w-[16rem] rounded-full bg-sky-400/10 blur-3xl" />
        </div>
        <div className="mx-auto grid max-w-7xl items-center gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:gap-10">
          <div>
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-cyan-400/20 bg-cyan-400/10 px-4 py-2 text-sm text-cyan-200">
              Connected website + Go gateway
            </div>
            <div className="mb-6"><Brand /></div>
            <h1 className="display-font max-w-4xl text-4xl font-black tracking-tight text-white sm:text-5xl lg:text-7xl">
              Start Hosting Today <span className="bg-gradient-to-r from-cyan-300 via-sky-300 to-blue-400 bg-clip-text text-transparent">with Easy2host.</span>
            </h1>
            <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300 sm:mt-6 sm:text-lg sm:leading-8">
              Create your account, select a plan, and launch your Minecraft server instantly. Manage everything from your personal dashboard with full control and flexibility.
            </p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:gap-4">
              <button onClick={() => setScreen(currentUser ? currentUser.role === "admin" ? "admin" : "dashboard" : "signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-6 py-3 font-semibold text-slate-950", popClass())}>
                {currentUser ? "Open Dashboard" : "Create Account"}
              </button>
              <a href="#plans" className={cn("rounded-2xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-white", popClass())}>Explore Plans</a>
            </div>
            <div className="mt-10 grid max-w-2xl gap-4 sm:grid-cols-3">
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">1.8 → 1.21.2</div><div className="mt-1 text-sm text-slate-300">Version support</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">Dashboards</div><div className="mt-1 text-sm text-slate-300">Clean & simple</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">Backups</div><div className="mt-1 text-sm text-slate-300">Create & download</div></GlassCard>
            </div>
          </div>

          <GlassCard className={cn("overflow-hidden p-6 lg:p-7", popClass())}>
            <div className="mb-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-400">Getting started</p>
                <h3 className="display-font text-xl font-semibold text-white">Four easy steps</h3>
              </div>
              <div className="rounded-xl border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-sm text-cyan-200">4 steps</div>
            </div>
            <div className="space-y-4">
              {[
                ["1", "Create an account", "Sign up with a username and password — takes about 30 seconds."],
                ["2", "Pick a plan", "Choose Starter, Standard, or Premium and add any extras you want."],
                ["3", "Launch your server", "We provision your Minecraft server instantly after checkout."],
                ["4", "You're in control", "Start, stop, back up, and monitor your server from your dashboard."],
              ].map(([number, title, text]) => (
                <div key={number} className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}>
                  <div className="flex items-start gap-4">
                    <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">{number}</div>
                    <div>
                      <div className="text-lg font-semibold text-white">{title}</div>
                      <div className="mt-1 text-sm text-slate-400">{text}</div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </GlassCard>
        </div>
      </section>

      <section id="plans" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <SectionIntro eyebrow="Plans" title="Pick the right plan for your server" text="All plans include a dedicated Minecraft server, full dashboard access, and one-click start/stop controls. Upgrade anytime." />
        <div className="mt-10 grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
          {plans.map((plan) => {
            const features = [
              `${plan.ram} GB RAM`,
              `${plan.storage} GB SSD storage`,
              `${plan.cores} CPU cores`,
              plan.backups === 0 ? "No included backups" : `${plan.backups} backup slot${plan.backups > 1 ? "s" : ""} included`,
              "Full dashboard control",
              "Start / stop anytime",
            ];
            return (
              <button
                key={plan.name}
                type="button"
                onClick={() => startPurchase(plan)}
                className={cn(
                  "relative isolate flex h-full w-full flex-col rounded-[2rem] border p-6 text-left",
                  popClass(),
                  plan.featured
                    ? "border-cyan-300/30 bg-gradient-to-b from-cyan-400/10 to-transparent"
                    : "border-white/10 bg-white/5",
                )}
              >
                {plan.featured && (
                  <div className="pointer-events-none absolute right-5 top-5 select-none rounded-full bg-gradient-to-r from-cyan-300 to-blue-500 px-3 py-1 text-xs font-bold uppercase tracking-[0.2em] text-slate-950">
                    Most popular
                  </div>
                )}
                <div>
                  <h3 className="display-font text-2xl font-bold text-white">{plan.name}</h3>
                  <div className="mt-2 text-sm text-slate-400">Perfect for {plan.name === "Starter" ? "small friend groups" : plan.name === "Standard" ? "growing communities" : "large servers"}</div>
                </div>
                <ul className="mt-6 flex-1 space-y-2.5">
                  {features.map((feat) => (
                    <li key={feat} className="flex items-center gap-2.5 text-sm text-slate-300">
                      <Check className="h-4 w-4 shrink-0 text-cyan-300" />
                      {feat}
                    </li>
                  ))}
                </ul>
                <div className="mt-8">
                  <div className="display-font text-3xl font-black text-white">
                    {money(plan.price)}
                    <span className="text-base font-medium text-slate-400"> / month</span>
                  </div>
                  <div className="mt-4 flex items-center gap-2 rounded-2xl border border-white/15 bg-white/8 px-4 py-3 font-semibold text-white">
                    Select plan <ChevronRight className="h-4 w-4" />
                  </div>
                </div>
              </button>
            );
          })}
        </div>

        <div className="mt-12">
          <div className="mb-5">
            <h3 className="display-font text-2xl font-bold text-white">Optional upgrades</h3>
            <p className="mt-2 text-slate-300">Add extra resources during checkout. Applied to your server at provisioning time.</p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {addOns.map((addon) => (
              <GlassCard key={addon.key} className={cn("p-5", popClass())}>
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <div className="text-lg font-semibold text-white">{addon.label}</div>
                    <div className="mt-1 text-sm text-slate-400">Add-on, billed monthly</div>
                  </div>
                  <div className="rounded-2xl bg-cyan-400/10 px-3 py-2 text-sm font-semibold text-cyan-300">{money(addon.price)}</div>
                </div>
                <div className="mt-4 space-y-1.5 text-sm text-slate-300">
                  {addon.ram > 0 && <div className="flex items-center gap-2"><Check className="h-3.5 w-3.5 text-cyan-300" /> +{addon.ram} GB RAM</div>}
                  {addon.storage > 0 && <div className="flex items-center gap-2"><Check className="h-3.5 w-3.5 text-cyan-300" /> +{addon.storage} GB SSD</div>}
                  <div className="flex items-center gap-2"><Check className="h-3.5 w-3.5 text-cyan-300" /> Selectable during checkout</div>
                </div>
              </GlassCard>
            ))}
          </div>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <SectionIntro eyebrow="Features" title="Everything you need to manage your server cleanly" text="A complete control panel for your Minecraft server with all the tools you need to run a successful server." />
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {featureCards.map(([Icon, title, text]) => (
            <GlassCard key={title} className={cn("p-6", popClass())}>
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">
                <Icon className="h-6 w-6" />
              </div>
              <h3 className="display-font mt-5 text-xl font-semibold text-white">{title}</h3>
              <p className="mt-3 leading-7 text-slate-300">{text}</p>
            </GlassCard>
          ))}
        </div>
      </section>

      <section id="faq" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <SectionIntro eyebrow="FAQ" title="Frequently asked questions" text="Find answers to common questions about easy2host and how to get the most out of your server." />
        <div className="mt-10 grid gap-4">
          {faqItems.map(([question, answer]) => (
            <GlassCard key={question} className={cn("p-6", popClass())}>
              <h3 className="text-lg font-semibold text-white">{question}</h3>
              <p className="mt-2 text-slate-300">{answer}</p>
            </GlassCard>
          ))}
        </div>
      </section>

      <section className="px-4 pb-20 sm:px-6 sm:pb-24 lg:px-16">
        <div className="mx-auto max-w-7xl">
          <GlassCard className="overflow-hidden p-8 sm:p-10">
            <div className="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <p className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">Ready to get started?</p>
                <h2 className="display-font mt-3 text-3xl font-bold text-white sm:text-4xl">Create your account and launch your server today</h2>
                <p className="mt-4 max-w-2xl text-slate-300">Choose your plan, select your Minecraft version, and take full control of your server.</p>
              </div>
              <button onClick={() => setScreen(currentUser ? currentUser.role === "admin" ? "admin" : "dashboard" : "signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-6 py-3 font-semibold text-slate-950", popClass())}>
                {currentUser ? "Open Dashboard" : "Create Account"}
              </button>
            </div>
          </GlassCard>
        </div>
      </section>
    </>
  );
}

export default function App({ initialScreen = "landing" }) {
  const [screen, setScreen] = useState(initialScreen);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [token, setToken] = useState(() => localStorage.getItem("easy2host_token") || "");
  const [currentUser, setCurrentUser] = useState(() => safeJsonParse(localStorage.getItem("easy2host_user"), null));
  const [servers, setServers] = useState(() => safeJsonParse(localStorage.getItem("easy2host_servers"), []));
  const [profiles, setProfiles] = useState(() => safeJsonParse(localStorage.getItem("easy2host_profiles"), {}));
  const [tickets, setTickets] = useState(() => safeJsonParse(localStorage.getItem("easy2host_tickets"), []));
  const [authBusy, setAuthBusy] = useState(false);
  const [authError, setAuthError] = useState("");
  const [purchasePlan, setPurchasePlan] = useState(null);
  const [notices, setNotices] = useState([]);
  const [apiHealthy, setApiHealthy] = useState(false);

  useEffect(() => {
    localStorage.setItem("easy2host_servers", JSON.stringify(servers));
  }, [servers]);

  useEffect(() => {
    if (currentUser) localStorage.setItem("easy2host_user", JSON.stringify(currentUser));
    else localStorage.removeItem("easy2host_user");
  }, [currentUser]);

  useEffect(() => {
    if (token) localStorage.setItem("easy2host_token", token);
    else localStorage.removeItem("easy2host_token");
  }, [token]);

  useEffect(() => {
    localStorage.setItem("easy2host_profiles", JSON.stringify(profiles));
  }, [profiles]);

  useEffect(() => {
    localStorage.setItem("easy2host_tickets", JSON.stringify(tickets));
  }, [tickets]);

  useEffect(() => {
    let link = document.querySelector("link[rel='icon']");
    if (!link) {
      link = document.createElement("link");
      link.setAttribute("rel", "icon");
      document.head.appendChild(link);
    }
    link.setAttribute("type", "image/png");
    link.setAttribute("href", logoIcon);
  }, []);

  useEffect(() => {
    setScreen(initialScreen);
  }, [initialScreen]);

  useEffect(() => {
    if (screen === "dashboard" && !currentUser) {
      setScreen("signin");
      return;
    }
    if (screen === "admin" && !currentUser) {
      setScreen("signin");
      return;
    }
    if (screen === "admin" && currentUser?.role !== "admin") {
      setScreen("dashboard");
    }
  }, [screen, currentUser]);

  useEffect(() => {
    const onUnauthorized = (event) => {
      const message = event?.detail?.message || "Session expired. Please sign in again.";
      setToken("");
      setCurrentUser(null);
      setScreen("signin");
      setNotices((previous) => [{ id: Date.now() + Math.random(), type: "error", text: message }, ...previous].slice(0, 5));
    };

    window.addEventListener("easy2host:unauthorized", onUnauthorized);
    return () => window.removeEventListener("easy2host:unauthorized", onUnauthorized);
  }, []);

  useEffect(() => {
    let ignore = false;
    const checkHealth = async () => {
      try {
        const response = await fetch(`${API_BASE}/gateway/health`);
        if (!ignore) setApiHealthy(response.ok);
      } catch {
        if (!ignore) setApiHealthy(false);
      }
    };

    checkHealth();
    const timer = window.setInterval(checkHealth, 30000);
    return () => {
      ignore = true;
      window.clearInterval(timer);
    };
  }, []);

  const pushNotice = (type, text) => setNotices((previous) => [{ id: Date.now() + Math.random(), type, text }, ...previous].slice(0, 5));

  const startPurchase = (plan) => setPurchasePlan(plan);

  const openLandingSection = (sectionId) => {
    const targetUrl = `/index.html#${sectionId}`;
    if (window.location.pathname !== "/index.html") {
      window.location.href = targetUrl;
      return;
    }
    const target = document.getElementById(sectionId);
    if (target) target.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  const logout = () => {
    setToken("");
    setCurrentUser(null);
    setScreen(initialScreen === "admin" || initialScreen === "dashboard" ? "signin" : initialScreen);
    setMobileOpen(false);
  };

  const buildUserFromToken = (username, jwt) => {
    const decoded = decodeJwt(jwt) || {};
    const isAdmin = Boolean(decoded.admin_id);
    const profile = profiles[username] || {};
    return {
      username,
      displayName: profile.displayName || username,
      email: profile.email || "",
      role: isAdmin ? "admin" : "user",
      ownerKey: decoded.admin_id || decoded.user_id || username,
      raw: decoded,
    };
  };

  const handleSignUp = async ({ displayName, username, email, password }) => {
    setAuthBusy(true);
    setAuthError("");
    try {
      await apiFetch("/auth/user/register", {
        method: "POST",
        body: {
          username,
          password,
          email,
          type: "user",
          jwt: "",
        },
      });
      setProfiles((previous) => ({
        ...previous,
        [username]: { displayName: displayName || username, email },
      }));
      setScreen("signin");
      pushNotice("success", "Account created. Sign in with your username and password.");
    } catch (error) {
      setAuthError(error.message);
    } finally {
      setAuthBusy(false);
    }
  };

  const handleSignIn = async ({ username, password }) => {
    setAuthBusy(true);
    setAuthError("");
    try {
      const data = await apiFetch("/auth/user/login", { method: "POST", body: { username, password } });
      const nextToken = data.token || "";
      const user = buildUserFromToken(username, nextToken);
      setToken(nextToken);
      setCurrentUser(user);
      setScreen(user.role === "admin" ? "admin" : "dashboard");
      pushNotice("success", `Signed in as ${user.username}.`);
    } catch (error) {
      setAuthError(error.message);
    } finally {
      setAuthBusy(false);
    }
  };

  const completePurchase = async (order, setCreateStep) => {
    if (!currentUser || !token) {
      setScreen("signin");
      return;
    }

    try {
      const bundleName = order.plan.name;
      pushNotice("info", `Creating your server \"${order.setup.name}\"...`);

      setCreateStep && setCreateStep("Step 1/3 — Generating bundle key…");
      const bundleKeyResponse = await apiFetch("/bundle/create", {
        method: "POST",
        body: { bundle: bundleName },
        token,
      });

      setCreateStep && setCreateStep("Step 2/3 — Adding bundle to your account…");
      await apiFetch("/bundle/add", { method: "POST", body: { bundle: bundleName }, token });

      setCreateStep && setCreateStep("Step 3/3 — Provisioning Minecraft server…");
      const created = await apiFetch("/server/create", {
        method: "POST",
        body: { version: order.setup.version, bundle: bundleKeyResponse.key, server_type: order.setup.software },
        token,
      });

      const extraRam = order.addons.reduce((sum, addon) => sum + addon.ram, 0);
      const extraStorage = order.addons.reduce((sum, addon) => sum + addon.storage, 0);
      const nextServer = {
        ownerId: currentUser.ownerKey,
        server_id: created.server_id,
        port: created.port || FIRST_SERVER_PORT + servers.length,
        name: order.setup.name,
        version: order.setup.version,
        software: order.setup.software,
        bundleName,
        status: "created",
        ram: order.plan.ram + extraRam,
        storage: order.plan.storage + extraStorage,
        cores: order.plan.cores,
        backupsIncluded: order.plan.backups,
        backupCount: 0,
      };
      setServers((previous) => [...previous, nextServer]);
      setPurchasePlan(null);
      setScreen("dashboard");
      pushNotice("success", `Server "${order.setup.name}" created — ID: ${created.server_id}.`);
    } catch (error) {
      pushNotice("error", `Server creation failed: ${error.message}`);
    }
  };

  const openPage = (nextScreen) => {
    const routes = {
      landing: "/index.html",
      signin: "/signin.html",
      signup: "/signup.html",
      dashboard: "/dashboard.html",
      admin: "/admin.html",
    };
    const target = routes[nextScreen] || "/index.html";
    if (window.location.pathname !== target) {
      window.location.href = target;
      return;
    }
    setScreen(nextScreen);
  };

  return (
    <div className="min-h-screen bg-[#07111f] text-white">
      <div className="fixed inset-0 -z-10 bg-[radial-gradient(circle_at_top,rgba(56,189,248,0.08),transparent_28%),linear-gradient(180deg,#08111f_0%,#09121c_45%,#070d17_100%)]" />

      <header className="sticky top-0 z-40 border-b border-white/10 bg-slate-950/70 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-4 sm:px-6 lg:px-16">
          <button onClick={() => openPage("landing")} className="min-w-0 text-left">
            <Brand />
          </button>

          <nav className="hidden items-center gap-8 md:flex">
            <button onClick={() => openPage("landing")} className="text-sm font-medium text-slate-300 transition hover:text-white">Home</button>
            <button onClick={() => openLandingSection("plans")} className="text-sm font-medium text-slate-300 transition hover:text-white">Plans</button>
            <button onClick={() => openLandingSection("features")} className="text-sm font-medium text-slate-300 transition hover:text-white">Features</button>
            <button onClick={() => openLandingSection("faq")} className="text-sm font-medium text-slate-300 transition hover:text-white">FAQ</button>
            {currentUser?.role === "user" && <button onClick={() => openPage("dashboard")} className="text-sm font-medium text-slate-300 transition hover:text-white">Server Dashboard</button>}
            {currentUser?.role === "admin" && <button onClick={() => openPage("admin")} className="text-sm font-medium text-slate-300 transition hover:text-white">Admin Dashboard</button>}
          </nav>

          <div className="hidden items-center gap-2 md:flex lg:gap-3">
            {currentUser ? (
              <>
                <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-slate-300">{currentUser.username}</div>
                <button onClick={() => openPage(currentUser.role === "admin" ? "admin" : "dashboard")} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 font-medium text-white", popClass())}>{currentUser.role === "admin" ? "Admin Panel" : "Dashboard"}</button>
              </>
            ) : (
              <>
                <button onClick={() => openPage("signin")} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 font-medium text-white", popClass())}>Sign In</button>
                <button onClick={() => openPage("signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-4 py-2.5 font-semibold text-slate-950", popClass())}>Create Account</button>
              </>
            )}
          </div>

          <button onClick={() => setMobileOpen((value) => !value)} className="rounded-2xl border border-white/10 bg-white/5 p-2.5 text-white md:hidden">
            {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </button>
        </div>

        {mobileOpen && (
          <div className="border-t border-white/10 px-4 py-4 sm:px-6 md:hidden">
            <div className="flex flex-col gap-3">
              <button onClick={() => { openPage("landing"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Home</button>
              <button onClick={() => { openLandingSection("plans"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Plans</button>
              <button onClick={() => { openLandingSection("features"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Features</button>
              <button onClick={() => { openLandingSection("faq"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">FAQ</button>
              {!currentUser && (
                <>
                  <button onClick={() => { openPage("signin"); setMobileOpen(false); }} className="rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-left text-white">Sign In</button>
                  <button onClick={() => { openPage("signup"); setMobileOpen(false); }} className="rounded-xl bg-gradient-to-r from-cyan-300 to-blue-500 px-3 py-2 text-left font-semibold text-slate-950">Create Account</button>
                </>
              )}
              {currentUser?.role === "user" && <button onClick={() => { openPage("dashboard"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Server Dashboard</button>}
              {currentUser?.role === "admin" && <button onClick={() => { openPage("admin"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Admin Dashboard</button>}
            </div>
          </div>
        )}
      </header>

      {screen === "landing" && <LandingPage apiHealthy={apiHealthy} currentUser={currentUser} setScreen={setScreen} startPurchase={startPurchase} />}
      {screen === "signup" && <AuthScreen mode="signup" busy={authBusy} error={authError} onSubmit={handleSignUp} setScreen={setScreen} />}
      {screen === "signin" && <AuthScreen mode="signin" busy={authBusy} error={authError} onSubmit={handleSignIn} setScreen={setScreen} />}
      {screen === "dashboard" && currentUser?.role === "user" && (
        <CustomerDashboard
          currentUser={currentUser}
          token={token}
          servers={servers}
          setServers={setServers}
          notices={notices}
          setNotices={setNotices}
          apiHealthy={apiHealthy}
          logout={logout}
          tickets={tickets}
          setTickets={setTickets}
          onBuyServer={startPurchase}
        />
      )}
      {screen === "admin" && currentUser?.role === "admin" && (
        <AdminDashboard
          currentUser={currentUser}
          token={token}
          notices={notices}
          setNotices={setNotices}
          apiHealthy={apiHealthy}
          logout={logout}
          tickets={tickets}
          setTickets={setTickets}
        />
      )}

      <PurchaseFlow open={Boolean(purchasePlan)} plan={purchasePlan} onClose={() => setPurchasePlan(null)} currentUser={currentUser} onRequireAuth={() => setScreen("signin")} onComplete={completePurchase} />
      <ToastStack notices={notices} setNotices={setNotices} />

      <footer className="border-t border-white/10 bg-slate-950/70 px-4 py-10 sm:px-6 lg:px-16">
        <div className="mx-auto grid max-w-7xl gap-6 md:grid-cols-3 md:items-end">
          <div>
            <div className="display-font text-xl font-bold text-white">easy2host</div>
            <p className="mt-2 max-w-md text-sm leading-6 text-slate-400">
              Host your Minecraft server with confidence. Easy setup, powerful controls, and reliable uptime—all in one place.
            </p>
          </div>
          <div className="text-sm text-slate-300 md:text-center">
            <div className="font-semibold text-white">Quick Links</div>
            <div className="mt-2 flex flex-wrap items-center gap-3 md:justify-center">
              <button type="button" onClick={() => openPage("landing")} className="text-slate-300 transition hover:text-cyan-200">Home</button>
              <button type="button" onClick={() => openLandingSection("plans")} className="text-slate-300 transition hover:text-cyan-200">Plans</button>
              <button type="button" onClick={() => openLandingSection("features")} className="text-slate-300 transition hover:text-cyan-200">Features</button>
              <button type="button" onClick={() => openLandingSection("faq")} className="text-slate-300 transition hover:text-cyan-200">FAQ</button>
            </div>
          </div>
          <div className="text-left text-sm text-slate-400 md:text-right">
            <div>Simple Minecraft hosting</div>
            <div className="mt-2">Copyright {new Date().getFullYear()} easy2host</div>
          </div>
        </div>
      </footer>
    </div>
  );
}

function UserSupportPanel({ currentUser, tickets, setTickets, pushNotice }) {
  const [form, setForm] = useState({ subject: "", message: "", priority: "normal" });

  const myTickets = useMemo(
    () => tickets.filter((ticket) => ticket.ownerId === currentUser.ownerKey),
    [tickets, currentUser.ownerKey],
  );

  const submitTicket = () => {
    const subject = form.subject.trim();
    const message = form.message.trim();
    if (!subject || !message) {
      pushNotice("error", "Please fill in both subject and message.");
      return;
    }

    const now = new Date().toISOString();
    const ticket = {
      id: makeTicketId(),
      ownerId: currentUser.ownerKey,
      ownerUsername: currentUser.username,
      subject,
      message,
      priority: form.priority,
      status: "open",
      adminReply: "",
      createdAt: now,
      updatedAt: now,
    };

    setTickets((previous) => [ticket, ...previous]);
    setForm({ subject: "", message: "", priority: "normal" });
    pushNotice("success", `Support ticket ${ticket.id} created.`);
  };

  const closeTicket = (ticketID) => {
    const now = new Date().toISOString();
    setTickets((previous) =>
      previous.map((ticket) =>
        ticket.id === ticketID
          ? {
              ...ticket,
              status: "resolved",
              updatedAt: now,
            }
          : ticket,
      ),
    );
    pushNotice("success", `Ticket ${ticketID} closed.`);
  };

  const deleteTicket = (ticketID) => {
    setTickets((previous) => previous.filter((ticket) => ticket.id !== ticketID));
    pushNotice("success", `Ticket ${ticketID} deleted.`);
  };

  return (
    <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
      <GlassCard className="p-6">
        <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Support</div>
        <h2 className="display-font mt-2 text-2xl font-bold text-white">Create Ticket</h2>
        <p className="mt-3 text-slate-300">Open a support ticket for billing, server setup, or runtime issues. Admins can handle it in the admin dashboard.</p>

        <div className="mt-6 space-y-4">
          <input
            value={form.subject}
            onChange={(event) => setForm((current) => ({ ...current, subject: event.target.value }))}
            placeholder="Subject"
            className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
          />

          <select
            value={form.priority}
            onChange={(event) => setForm((current) => ({ ...current, priority: event.target.value }))}
            className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none"
          >
            <option value="low">Low</option>
            <option value="normal">Normal</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </select>

          <textarea
            value={form.message}
            onChange={(event) => setForm((current) => ({ ...current, message: event.target.value }))}
            placeholder="Describe your issue"
            rows={6}
            className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
          />

          <div className="flex justify-end">
            <button
              type="button"
              onClick={submitTicket}
              className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950", popClass())}
            >
              Submit Ticket
            </button>
          </div>
        </div>
      </GlassCard>

      <GlassCard className="p-6">
        <div className="mb-4 flex items-center justify-between gap-3">
          <h3 className="display-font text-2xl font-bold text-white">My Tickets</h3>
          <div className="rounded-2xl bg-white/8 px-3 py-1 text-sm text-slate-300">{myTickets.length}</div>
        </div>

        <div className="space-y-3">
          {myTickets.length === 0 && (
            <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">
              No tickets yet.
            </div>
          )}

          {myTickets.map((ticket) => (
            <div key={ticket.id} className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="font-semibold text-white">{ticket.subject}</div>
                <div className="rounded-full bg-white/8 px-3 py-1 text-xs uppercase tracking-[0.12em] text-slate-300">{ticket.status}</div>
              </div>

              <div className="mt-2 text-xs uppercase tracking-[0.12em] text-slate-500">{ticket.id} · {ticket.priority}</div>
              <p className="mt-3 text-sm leading-6 text-slate-300">{ticket.message}</p>
              {ticket.adminReply && (
                <div className="mt-3 rounded-xl border border-cyan-300/20 bg-cyan-400/10 p-3 text-sm text-cyan-100">
                  <div className="font-semibold">Admin reply</div>
                  <div className="mt-1">{ticket.adminReply}</div>
                </div>
              )}
              {ticket.status !== "resolved" && (
                <div className="mt-3 flex justify-end gap-2">
                  <button
                    type="button"
                    onClick={() => closeTicket(ticket.id)}
                    className={cn("rounded-2xl border border-emerald-400/30 bg-emerald-400/10 px-4 py-2 text-sm font-semibold text-emerald-200", popClass())}
                  >
                    Close Ticket
                  </button>
                  <button
                    type="button"
                    onClick={() => deleteTicket(ticket.id)}
                    className={cn("rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-2 text-sm font-semibold text-rose-200", popClass())}
                  >
                    Delete Ticket
                  </button>
                </div>
              )}
              {ticket.status === "resolved" && (
                <div className="mt-3 flex justify-end">
                  <button
                    type="button"
                    onClick={() => deleteTicket(ticket.id)}
                    className={cn("rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-2 text-sm font-semibold text-rose-200", popClass())}
                  >
                    Delete Ticket
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </GlassCard>
    </div>
  );
}

function AdminTicketPanel({ tickets, setTickets, pushNotice }) {
  const [selectedId, setSelectedId] = useState("");
  const [reply, setReply] = useState("");
  const [status, setStatus] = useState("open");

  const selectedTicket = tickets.find((ticket) => ticket.id === selectedId) || null;

  // Use selectedId as dep to avoid object-reference churn causing stale reads
  useEffect(() => {
    const ticket = tickets.find((t) => t.id === selectedId) || null;
    if (!ticket) {
      setReply("");
      setStatus("open");
      return;
    }
    setReply(ticket.adminReply || "");
    setStatus(ticket.status || "open");
  }, [selectedId, tickets]);

  const saveTicket = () => {
    if (!selectedTicket) {
      pushNotice("error", "Select a ticket first.");
      return;
    }

    const now = new Date().toISOString();
    setTickets((previous) =>
      previous.map((ticket) =>
        ticket.id === selectedTicket.id
          ? {
              ...ticket,
              status,
              adminReply: reply.trim(),
              updatedAt: now,
            }
          : ticket,
      ),
    );

    pushNotice("success", `Ticket ${selectedTicket.id} updated.`);
  };

  return (
    <div className="grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
      <GlassCard className="p-6">
        <div className="mb-4 flex items-center justify-between gap-3">
          <h2 className="display-font text-2xl font-bold text-white">Support Queue</h2>
          <div className="rounded-2xl bg-white/8 px-3 py-1 text-sm text-slate-300">{tickets.length}</div>
        </div>
        <div className="space-y-3">
          {tickets.length === 0 && (
            <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">
              No open tickets.
            </div>
          )}
          {tickets.map((ticket) => (
            <button
              key={ticket.id}
              type="button"
              onClick={() => setSelectedId(ticket.id)}
              className={cn(
                "w-full rounded-2xl border p-4 text-left",
                popClass(),
                selectedId === ticket.id ? "border-cyan-300/30 bg-cyan-400/10" : "border-white/10 bg-slate-950/40",
              )}
            >
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="font-semibold text-white">{ticket.subject}</div>
                <div className="rounded-full bg-white/8 px-3 py-1 text-xs uppercase tracking-[0.12em] text-slate-300">{ticket.status}</div>
              </div>
              <div className="mt-2 text-xs uppercase tracking-[0.12em] text-slate-500">{ticket.id} · {ticket.ownerUsername} · {ticket.priority}</div>
            </button>
          ))}
        </div>
      </GlassCard>

      <GlassCard className="p-6">
        {!selectedTicket ? (
          <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">
            Select a ticket to update status and reply.
          </div>
        ) : (
          <>
            <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Ticket Details</div>
            <h3 className="display-font mt-2 text-2xl font-bold text-white">{selectedTicket.subject}</h3>
            <div className="mt-2 text-xs uppercase tracking-[0.12em] text-slate-500">{selectedTicket.id} · {selectedTicket.ownerUsername} · {selectedTicket.priority}</div>
            <p className="mt-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4 text-sm leading-6 text-slate-300">{selectedTicket.message}</p>

            <div className="mt-5 grid gap-4">
              <select
                value={status}
                onChange={(event) => setStatus(event.target.value)}
                className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none"
              >
                <option value="open">Open</option>
                <option value="in_progress">In Progress</option>
                <option value="resolved">Resolved</option>
              </select>

              <textarea
                value={reply}
                onChange={(event) => setReply(event.target.value)}
                rows={6}
                placeholder="Write an admin response"
                className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none"
              />

              <div className="flex justify-end">
                <button
                  type="button"
                  onClick={saveTicket}
                  className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950", popClass())}
                >
                  Save Ticket Update
                </button>
              </div>
            </div>
          </>
        )}
      </GlassCard>
    </div>
  );
}