import React, { useEffect, useMemo, useState } from "react";
import {
  Activity,
  Boxes,
  Check,
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
  Play,
  Plus,
  Search,
  Server,
  Settings,
  Shield,
  Square,
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
const softwareOptions = ["Vanilla", "Fabric", "Bukkit", "Paper"];

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
  return "anim-fade-up transition duration-200 motion-safe:hover:-translate-y-1 motion-safe:hover:scale-[1.02]";
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

function normalizeStatsValue(value) {
  if (typeof value !== "string") return value;
  const parsed = safeJsonParse(value, null);
  return parsed ?? value;
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
          <div className="text-[11px] uppercase tracking-[0.32em] text-slate-400">Minecraft control plane</div>
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

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center bg-black/75 p-2 sm:items-center sm:p-4"
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
    </div>
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

function NoticeStack({ notices }) {
  if (notices.length === 0) return null;

  return (
    <GlassCard className="p-4">
      {notices.map((notice) => (
        <div
          key={notice.id}
          className={cn(
            "mb-2 rounded-2xl px-4 py-3 text-sm last:mb-0",
            notice.type === "error" ? "bg-rose-400/10 text-rose-100" : "bg-cyan-400/10 text-cyan-100",
          )}
        >
          {notice.text}
        </div>
      ))}
    </GlassCard>
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
      return;
    }
    if (step === 2 && !billing.fullName.trim()) return;
    if (step === 3 && !setup.name.trim()) return;
    if (step < 3) {
      setStep((value) => value + 1);
      return;
    }
    setCreating(true);
    onComplete({ plan, addons: selected, setup, billing }, setCreateStep).finally(() => setCreating(false));
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={`Buy ${plan.name}`}
      subtitle="Choose extras, review the order, and send the provisioning request through the existing gateway API."
    >
      <div className="mb-6 flex flex-wrap gap-3">
        {["Extras", "Billing", "Server Setup"].map((label, index) => (
          <div
            key={label}
            className={cn(
              "rounded-full px-4 py-2 text-sm",
              step === index + 1 ? "bg-cyan-400/15 text-cyan-300" : "bg-white/5 text-slate-400",
            )}
          >
            {label}
          </div>
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
              Billing remains UI-only for now. The actual bundle key and server provisioning call happens on the final step.
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
          <select
            value={setup.version}
            onChange={(event) => setSetup((current) => ({ ...current, version: event.target.value }))}
            className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none"
          >
            {versions.map((version) => (
              <option key={version}>{version}</option>
            ))}
          </select>
          <div className="rounded-2xl border border-dashed border-cyan-300/25 bg-cyan-400/5 px-4 py-3 text-slate-300">
            The selected plan becomes the bundle key request. Software and version are stored locally for the dashboard and version provisioning call.
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
              <details className="rounded-2xl border border-amber-400/20 bg-amber-400/5 p-4 text-sm text-slate-300">
                <summary className="cursor-pointer font-semibold text-amber-200">How to create an admin account</summary>
                <div className="mt-3 space-y-2 text-slate-300">
                  <p>Admin accounts require knowing the <strong className="text-white">JWT_SECRET</strong> from your server config (<code className="rounded bg-slate-800 px-1">hub/config/docker.env</code>).</p>
                  <p>Run this command (replace <code className="rounded bg-slate-800 px-1">YOUR_JWT_SECRET</code>, choose your own username/password):</p>
                  <pre className="mt-2 overflow-auto rounded-xl bg-slate-900 p-3 text-xs text-cyan-200">{`curl -s -X POST http://localhost:8010/api/v1/auth/user/register \\
  -H "Content-Type: application/json" \\
  -d '{"username":"admin","password":"yourpassword","email":"","type":"admin","jwt":"YOUR_JWT_SECRET"}'`}</pre>
                  <p className="text-slate-400">Then sign in normally. Admin credentials automatically open the Admin Dashboard.</p>
                </div>
              </details>
            </div>
          )}
        </GlassCard>
        <GlassCard className="hidden overflow-hidden lg:block">
          <div className="h-full bg-[radial-gradient(circle_at_top_left,rgba(103,232,249,0.16),transparent_38%),linear-gradient(160deg,rgba(8,17,31,0.94),rgba(10,20,34,0.92))] p-8">
            <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-cyan-400/20 bg-cyan-400/10 px-4 py-2 text-sm text-cyan-200">
              Connected to your Go gateway
            </div>
            <div className="space-y-5">
              {[
                "Auth uses /api/v1/auth/user/register and /api/v1/auth/user/login.",
                "Admin account: register with type=admin and jwt=<JWT secret>, then login normally.",
                "Purchase flow calls /bundle/create, /bundle/add, then /server/create.",
                "Customer dashboard controls start, stop, stats, delete, and backup download.",
                "Admin dashboard reaches network and host metadata routes through the same gateway.",
              ].map((line) => (
                <div key={line} className="flex items-start gap-3 rounded-2xl border border-white/10 bg-white/5 p-4 text-slate-200">
                  <Check className="mt-0.5 h-5 w-5 shrink-0 text-cyan-300" />
                  <span>{line}</span>
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
          <p className="mt-2 text-slate-400">This code remains in place but is hidden from the visible product flow for now.</p>
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

function CustomerDashboard({ currentUser, token, servers, setServers, notices, setNotices, apiHealthy, logout }) {
  const [active, setActive] = useState("servers");
  const [selectedServerId, setSelectedServerId] = useState(null);
  const [search, setSearch] = useState("");
  const [statsMap, setStatsMap] = useState({});
  const [serverBusy, setServerBusy] = useState({});

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
  const pushNotice = (type, text) => setNotices((previous) => [{ id: Date.now() + Math.random(), type, text }, ...previous].slice(0, 5));

  const updateServer = (serverId, updater) => {
    setServers((previous) => previous.map((server) => (server.server_id === serverId ? { ...server, ...updater(server) } : server)));
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

  const fetchStats = async (server) => {
    setBusy(server.server_id, true);
    try {
      const data = await apiFetch("/server/getStats", {
        method: "POST",
        body: { server_id: server.server_id, key: "Online" },
        token,
      });
      setStatsMap((previous) => ({ ...previous, [server.server_id]: normalizeStatsValue(data.value) }));
      pushNotice("success", `Loaded stats for ${server.name}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  useEffect(() => {
    if (!selectedServer || !token) return;

    const timer = window.setInterval(() => {
      fetchStats(selectedServer);
    }, 30000);

    return () => window.clearInterval(timer);
  }, [selectedServer, token]);

  const createBackup = async (server) => {
    setBusy(server.server_id, true);
    try {
      await apiFetch("/server/backup/create", {
        method: "POST",
        body: { server_id: server.server_id, bundle: server.bundleName },
        token,
      });
      updateServer(server.server_id, (current) => ({ backupCount: (current.backupCount || 0) + 1 }));
      pushNotice("success", `Backup created for ${server.name}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const downloadBackup = async (server) => {
    setBusy(server.server_id, true);
    try {
      const blob = await apiFetchBlob("/server/backup/get", { method: "POST", body: { server_id: server.server_id }, token });
      downloadBlob(blob, `${server.name || server.server_id}-backup.tar.gz`);
      pushNotice("success", `Downloaded backup for ${server.name}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const menu = [
    { key: "servers", label: "My Servers", icon: Server },
    ...(SHOW_PLUGIN_CATALOG ? [{ key: "catalog", label: "Plugins & Mods", icon: Boxes }] : []),
    { key: "account", label: "Account", icon: Settings },
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-10 lg:px-16">
      <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">Server Dashboard</div>
          <h1 className="display-font mt-2 text-3xl font-bold text-white">Welcome back, {currentUser.displayName}</h1>
          <p className="mt-2 text-slate-300">This panel is connected to your auth, bundle, server, stats, and backup routes.</p>
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
        <Sidebar items={menu} active={active} setActive={setActive} title="easy2host Panel" subtitle="User API area" />
        <div className="space-y-6">
          <NoticeStack notices={notices} />

          {active === "servers" && (
            <>
              <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
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
                <GlassCard className={cn("p-5", popClass())}>
                  <div className="text-sm text-slate-400">Next host port hint</div>
                  <div className="display-font mt-2 text-2xl font-bold text-cyan-300">
                    {ownedServers.reduce((highest, server) => Math.max(highest, server.port || 0), FIRST_SERVER_PORT - 1) + 1}
                  </div>
                </GlassCard>
              </div>

              <div className="grid gap-6 2xl:grid-cols-[0.95fr_1.05fr]">
                <GlassCard className="p-5">
                  <div className="mb-4 flex items-center justify-between gap-4">
                    <h2 className="display-font text-xl font-semibold text-white">My Servers</h2>
                    <div className="text-sm text-slate-400">Tracked locally until a backend list endpoint exists.</div>
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
                          "w-full rounded-2xl border p-4 text-left",
                          popClass(),
                          selectedServer?.server_id === server.server_id
                            ? "border-cyan-300/30 bg-cyan-400/10"
                            : "border-white/10 bg-slate-950/40 hover:bg-white/5",
                        )}
                      >
                        <div className="flex items-center justify-between gap-4">
                          <div>
                            <div className="font-semibold text-white">{server.name}</div>
                            <div className="mt-1 text-sm text-slate-400">
                              {server.software} {server.version} · {server.bundleName} · {server.server_id}
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
                          <p className="mt-1 break-all text-xs uppercase tracking-[0.18em] text-slate-500">{selectedServer.server_id}</p>
                        </div>
                        <div className={cn("rounded-full px-3 py-1 text-sm", selectedServer.status === "online" ? "bg-emerald-400/10 text-emerald-300" : "bg-slate-700/60 text-slate-300")}>
                          {selectedServer.status || "created"}
                        </div>
                      </div>

                      <div className="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "start")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-emerald-400/15 px-4 py-3 text-emerald-300 disabled:opacity-60", popClass())}><Play className="h-4 w-4" /> Start</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "stop")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-400/15 px-4 py-3 text-rose-300 disabled:opacity-60", popClass())}><Square className="h-4 w-4" /> Stop</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => fetchStats(selectedServer)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-sky-400/15 px-4 py-3 text-sky-300 disabled:opacity-60", popClass())}><Activity className="h-4 w-4" /> Stats</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => createBackup(selectedServer)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-violet-400/15 px-4 py-3 text-violet-300 disabled:opacity-60", popClass())}><HardDrive className="h-4 w-4" /> Create Backup</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => downloadBackup(selectedServer)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-white/8 px-4 py-3 text-white disabled:opacity-60", popClass())}><Download className="h-4 w-4" /> Download Backup</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer, "delete")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-500/20 px-4 py-3 text-rose-200 disabled:opacity-60", popClass())}><Trash2 className="h-4 w-4" /> Delete</button>
                      </div>

                      <div className="mt-6 grid gap-4 xl:grid-cols-2">
                        <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex items-center gap-2 text-white"><Database className="h-4 w-4 text-cyan-300" /> Server Overview</div>
                          <pre className="max-h-64 overflow-auto text-xs text-slate-300">{JSON.stringify({
                            server_id: selectedServer.server_id,
                            host_port: selectedServer.port,
                            plan: selectedServer.bundleName,
                            ram_gb: selectedServer.ram,
                            storage_gb: selectedServer.storage,
                            cpu_cores: selectedServer.cores,
                            status: selectedServer.status,
                            software: selectedServer.software,
                            version: selectedServer.version,
                            backups_created: selectedServer.backupCount || 0,
                          }, null, 2)}</pre>
                        </div>
                        <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex items-center gap-2 text-white"><Activity className="h-4 w-4 text-cyan-300" /> Stats Response</div>
                          <pre className="max-h-64 overflow-auto text-xs text-slate-300">{JSON.stringify(statsMap[selectedServer.server_id] || { message: "Load stats to see backend response." }, null, 2)}</pre>
                        </div>
                      </div>
                    </>
                  ) : (
                    <div className="text-slate-300">No server selected.</div>
                  )}
                </GlassCard>
              </div>
            </>
          )}

          {active === "catalog" && SHOW_PLUGIN_CATALOG && <PluginCatalogPanel search={search} setSearch={setSearch} />}

          {active === "account" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Account</div>
              <h2 className="display-font mt-2 text-2xl font-bold text-white">Profile & Session</h2>
              <div className="mt-8 grid gap-4 md:grid-cols-4">
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Display name</div><div className="mt-2 font-semibold text-white">{currentUser.displayName}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Username</div><div className="mt-2 font-semibold text-white">{currentUser.username}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Email</div><div className="mt-2 font-semibold text-white">{currentUser.email || "Stored locally"}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Role</div><div className="mt-2 inline-flex items-center gap-2 font-semibold text-white"><Lock className="h-4 w-4 text-cyan-300" /> {currentUser.role}</div></div>
              </div>
              <div className="mt-6 rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                <div className="mb-3 flex items-center gap-2 text-white"><KeyRound className="h-4 w-4 text-cyan-300" /> JWT</div>
                <pre className="overflow-auto text-xs text-slate-300">{token || "No token loaded."}</pre>
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

function AdminDashboard({ currentUser, token, notices, setNotices, apiHealthy, logout }) {
  const [active, setActive] = useState("overview");
  const [networkIp, setNetworkIp] = useState("");
  const [networkHostId, setNetworkHostId] = useState("");
  const [networkResult, setNetworkResult] = useState(null);
  const [hostMetadata, setHostMetadata] = useState(null);
  const [hostCreateForm, setHostCreateForm] = useState({ ram: "", cores: "" });
  const [hostDeleteId, setHostDeleteId] = useState("");
  const [hostAddForm, setHostAddForm] = useState({ host_server_id: "", server_id: "" });
  const [busy, setBusy] = useState(false);

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
          <p className="mt-2 text-slate-300">This uses your existing admin-guarded gateway endpoints rather than mocked admin data.</p>
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
          <NoticeStack notices={notices} />

          {active === "overview" && (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">API Base</div><div className="mt-2 text-lg font-bold text-cyan-300">{API_BASE}</div></GlassCard>
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Role</div><div className="display-font mt-2 text-3xl font-bold text-white">Admin</div></GlassCard>
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">JWT present</div><div className="display-font mt-2 text-3xl font-bold text-white">{token ? "Yes" : "No"}</div></GlassCard>
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Identity</div><div className="mt-2 text-sm font-semibold text-white">{currentUser.ownerKey}</div></GlassCard>
            </div>
          )}

          {active === "network" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Admin</div>
              <h2 className="display-font mt-2 text-2xl font-bold text-white">Create Network Host Entry</h2>
              <p className="mt-3 text-slate-300">Calls POST /network/create. For provisioning, use the host ID created in Host Metadata so IDs stay paired.</p>
              <div className="mt-6 grid gap-4 sm:grid-cols-2">
                <input value={networkIp} onChange={(event) => setNetworkIp(event.target.value)} placeholder="Host IP" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                <input value={networkHostId} onChange={(event) => setNetworkHostId(event.target.value)} placeholder="Host server ID (from metadata create)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
              </div>
              <div className="mt-4 flex">
                <button disabled={busy} onClick={async () => {
                  const payload = networkHostId.trim() ? { ip: networkIp, host_server_id: networkHostId.trim() } : { ip: networkIp };
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
                  <input value={hostCreateForm.ram} onChange={(event) => setHostCreateForm((current) => ({ ...current, ram: event.target.value }))} placeholder="RAM" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostCreateForm.cores} onChange={(event) => setHostCreateForm((current) => ({ ...current, cores: event.target.value }))} placeholder="CPU cores" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostDeleteId} onChange={(event) => setHostDeleteId(event.target.value)} placeholder="Host server ID to delete" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none md:col-span-2" />
                </div>
                <div className="mt-6 flex flex-wrap gap-3">
                  <button disabled={busy} onClick={() => callAdmin("Metadata create", () => apiFetch("/host-metadata/create", { method: "POST", body: hostCreateForm, token }))} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-4 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Create Metadata</button>
                  <button disabled={busy} onClick={async () => {
                    const result = await callAdmin("Metadata get", () => apiFetch("/host-metadata/get", { method: "GET", token }));
                    if (result) setHostMetadata(result);
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
                <input value={hostAddForm.host_server_id} onChange={(event) => setHostAddForm((current) => ({ ...current, host_server_id: event.target.value }))} placeholder="Host server ID" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                <input value={hostAddForm.server_id} onChange={(event) => setHostAddForm((current) => ({ ...current, server_id: event.target.value }))} placeholder="Server ID" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
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
        [Settings, "Version and software choice", "Pick Vanilla, Fabric, Bukkit, or Paper with the Minecraft version you want."],
        [CreditCard, "Simple purchase flow", "Bundle, billing, and provisioning move through one sequence."],
      ]
    : [
        [LayoutDashboard, "Clean dashboard", "One place for status, actions, backups, and account access."],
        [Server, "Easy server control", "Start, stop, back up, and inspect your server from one page."],
        [HardDrive, "Backup workflow", "Create a backup and download the tarball directly from the dashboard."],
        [Shield, "Account-based access", "Each account only sees and manages its own infrastructure."],
        [Settings, "Version and software choice", "Pick Vanilla, Fabric, Bukkit, or Paper with the Minecraft version you want."],
        [Crown, "Admin oversight", "The admin dashboard is available through the same sign-in flow and gateway."],
      ];

  const faqItems = [
    ["How fast can I create a server?", "Once your account is ready, you can choose a plan, finish checkout, and provision through the existing backend in a few steps."],
    ["Which server software can I choose?", "You can choose Vanilla, Fabric, Bukkit, or Paper during setup while still selecting the exact Minecraft version."],
    ["Can I add extra resources?", "Yes. During checkout you can add extra RAM and SSD upgrades before the server is created."],
    ["Will other users be able to see my server?", "No. Your dashboard only tracks the servers associated with your signed-in account."],
    ["Is there an admin panel?", "Yes. Admin credentials automatically open the admin dashboard and reach the protected infrastructure routes."],
    ["Are plugins gone?", "The catalog code is still preserved, but the visible plugin feature is intentionally hidden for now."],
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
              Sign up, choose a plan, generate a real bundle key, provision a server, and manage it through the server dashboard. Admin credentials land in the infrastructure dashboard automatically.
            </p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:gap-4">
              <button onClick={() => setScreen(currentUser ? currentUser.role === "admin" ? "admin" : "dashboard" : "signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-6 py-3 font-semibold text-slate-950", popClass())}>
                {currentUser ? "Open Dashboard" : "Create Account"}
              </button>
              <a href="#plans" className={cn("rounded-2xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-white", popClass())}>Explore Plans</a>
            </div>
            <div className="mt-10 grid max-w-2xl gap-4 sm:grid-cols-4">
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">1.8 → 1.21.2</div><div className="mt-1 text-sm text-slate-300">Version support</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">Dashboards</div><div className="mt-1 text-sm text-slate-300">User + admin</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">Backups</div><div className="mt-1 text-sm text-slate-300">Create + download</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="display-font text-2xl font-bold text-white">{apiHealthy ? "Online" : "Pending"}</div><div className="mt-1 text-sm text-slate-300">Gateway connection</div></GlassCard>
            </div>
          </div>

          <GlassCard className={cn("overflow-hidden p-6 lg:p-7", popClass())}>
            <div className="mb-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-400">Provisioning sequence</p>
                <h3 className="display-font text-xl font-semibold text-white">Real API flow</h3>
              </div>
              <div className="rounded-xl border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-sm text-cyan-200">4 steps</div>
            </div>
            <div className="space-y-4">
              {[
                ["1", "Authenticate", "Register or sign in through the auth service."],
                ["2", "Bundle key", "Generate a bundle key and add the selected plan to the user account."],
                ["3", "Provision server", "Create the Minecraft server with version, host selection, and tracked port."],
                ["4", "Operate it", "Use the server dashboard for start, stop, stats, backup, and delete."],
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
        <SectionIntro eyebrow="Plans" title="Choose a plan, then create a server" text="The visible purchase flow now maps to the existing bundle and task-service APIs instead of staying as a mock landing page." />
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {plans.map((plan) => (
            <button key={plan.name} onClick={() => startPurchase(plan)} className={cn("relative rounded-[2rem] border p-6 text-left", popClass(), plan.featured ? "border-cyan-300/30 bg-cyan-400/10 shadow-cyan-500/10" : "border-white/10 bg-white/5")}>
              {plan.featured && <div className="pointer-events-none absolute right-5 top-5 select-none rounded-full bg-gradient-to-r from-cyan-300 to-blue-500 px-3 py-1 text-xs font-bold uppercase tracking-[0.2em] text-slate-950">Most popular</div>}
              <h3 className="display-font text-2xl font-bold text-white">{plan.name}</h3>
              <div className="mt-6 space-y-2 text-slate-300">
                <div>{plan.ram} GB RAM</div>
                <div>{plan.storage} GB SSD</div>
                <div>{plan.cores} Cores</div>
                <div>{plan.backups} Backup slots</div>
              </div>
              <div className="display-font mt-8 text-3xl font-black text-white">{money(plan.price)} <span className="text-base font-medium text-slate-400">/ month</span></div>
              <div className="mt-8 inline-flex items-center gap-2 rounded-2xl bg-white/8 px-4 py-3 font-semibold text-white">Select plan <ChevronRight className="h-4 w-4" /></div>
            </button>
          ))}
        </div>

        <div className="mt-12">
          <div className="mb-5 flex items-center justify-between gap-4">
            <div>
              <h3 className="display-font text-2xl font-bold text-white">Available upgrades</h3>
              <p className="mt-2 text-slate-300">These upgrades are applied locally in the dashboard metadata so users can see the resources they bought.</p>
            </div>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {addOns.map((addon) => (
              <GlassCard key={addon.key} className={cn("p-5", popClass())}>
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <div className="text-lg font-semibold text-white">{addon.label}</div>
                    <div className="mt-2 text-sm text-slate-400">Monthly add-on for extra server resources.</div>
                  </div>
                  <div className="rounded-2xl bg-cyan-400/10 px-3 py-2 text-sm font-semibold text-cyan-300">{money(addon.price)}</div>
                </div>
                <div className="mt-5 space-y-2 text-sm text-slate-300">
                  {addon.ram > 0 && <div>+{addon.ram} GB RAM</div>}
                  {addon.storage > 0 && <div>+{addon.storage} GB SSD</div>}
                  <div>Selectable during checkout</div>
                </div>
              </GlassCard>
            ))}
          </div>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <SectionIntro eyebrow="Features" title="Everything you need to manage your server cleanly" text="The served website now reflects the real backend capabilities that exist today, and hides the plugin marketplace until that flow is ready." />
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
        <SectionIntro eyebrow="FAQ" title="Frequently asked questions" text="The frontend now mirrors the real contract and the remaining hidden feature surface is intentional." />
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
                <p className="mt-4 max-w-2xl text-slate-300">Pick a plan, configure the version, and manage everything from the connected dashboard.</p>
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
        body: { version: order.setup.version, bundle: bundleKeyResponse.key },
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
      {screen === "dashboard" && currentUser?.role === "user" && <CustomerDashboard currentUser={currentUser} token={token} servers={servers} setServers={setServers} notices={notices} setNotices={setNotices} apiHealthy={apiHealthy} logout={logout} />}
      {screen === "admin" && currentUser?.role === "admin" && <AdminDashboard currentUser={currentUser} token={token} notices={notices} setNotices={setNotices} apiHealthy={apiHealthy} logout={logout} />}

      <PurchaseFlow open={Boolean(purchasePlan)} plan={purchasePlan} onClose={() => setPurchasePlan(null)} currentUser={currentUser} onRequireAuth={() => setScreen("signin")} onComplete={completePurchase} />

      <footer className="border-t border-white/10 px-4 py-8 text-center text-sm text-slate-400 sm:px-6 lg:px-16">
        easy2host web app with connected auth, bundle, server lifecycle, backup, and admin infrastructure flows.
      </footer>
    </div>
  );
}