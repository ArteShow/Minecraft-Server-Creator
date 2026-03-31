import React, { useEffect, useMemo, useState } from "react";
import {
  Server,
  Shield,
  Play,
  Square,
  RotateCw,
  Search,
  LogIn,
  UserPlus,
  LayoutDashboard,
  Settings,
  Boxes,
  Menu,
  X,
  Check,
  CreditCard,
  Plus,
  Database,
  HardDrive,
  Upload,
  RefreshCw,
  Trash2,
  KeyRound,
  Activity,
  Lock,
} from "lucide-react";

const API_BASE = "http://localhost:8010/api/v1";
const ADMIN_EMAILS = ["gravitycrafter464@gmail.com"];

const plans = [
  { name: "Starter", ram: 2, cores: 2, cpuPower: 200, backups: 0, storage: 5, price: 5.99, featured: false },
  { name: "Standard", ram: 4, cores: 3, cpuPower: 300, backups: 1, storage: 15, price: 8.99, featured: true },
  { name: "Premium", ram: 5, cores: 4, cpuPower: 400, backups: 2, storage: 20, price: 10.99, featured: false },
];

const addOns = [
  { key: "ram1", label: "+1 GB RAM", price: 1.99, ram: 1, storage: 0 },
  { key: "ssd5", label: "+5 GB SSD", price: 1.49, ram: 0, storage: 5 },
  { key: "ssd10", label: "+10 GB SSD", price: 2.0, ram: 0, storage: 10 },
];

const versions = ["1.8.9", "1.12.2", "1.16.5", "1.20.6", "1.21.1", "1.21.2"];
const softwareOptions = ["Vanilla", "Fabric", "Bukkit", "Paper"];

const fallbackMods = [
  { id: "sodium", name: "Sodium", source: "Modrinth", type: "Performance", description: "Performance-focused rendering optimization mod." },
  { id: "lithium", name: "Lithium", source: "Modrinth", type: "Optimization", description: "General server and game logic optimizations." },
  { id: "fabric-api", name: "Fabric API", source: "Modrinth", type: "Core", description: "Core hooks and shared API for Fabric mods." },
];

function cn(...classes) {
  return classes.filter(Boolean).join(" ");
}

function money(value) {
  return new Intl.NumberFormat("en-EN", { style: "currency", currency: "EUR" }).format(value);
}

function popClass() {
  return "transition duration-200 motion-safe:hover:-translate-y-1 motion-safe:hover:scale-[1.02]";
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
        .map((c) => `%${(`00${c.charCodeAt(0).toString(16)}`).slice(-2)}`)
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

  if (!response.ok) {
    const message =
      data?.message ||
      data?.error ||
      data?.raw ||
      `Request failed with status ${response.status}`;
    throw new Error(message);
  }

  return data;
}

function GlassCard({ className = "", children }) {
  return (
    <div
      className={cn(
        "rounded-3xl border border-white/10 bg-white/5 backdrop-blur-xl shadow-2xl shadow-black/20",
        className,
      )}
    >
      {children}
    </div>
  );
}

function LogoMark({ full = false }) {
  return (
    <div className="flex items-center gap-3">
      <div className="relative flex h-11 w-11 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-300 via-sky-400 to-blue-600 shadow-lg shadow-cyan-500/20">
        <div className="absolute inset-1 rounded-xl bg-slate-950/20" />
        <Server className="relative z-10 h-5 w-5 text-white" />
      </div>
      <div>
        <div className="text-lg font-bold tracking-tight text-white">easy2host</div>
        {full && <div className="text-xs uppercase tracking-[0.25em] text-slate-400">simple minecraft hosting</div>}
      </div>
    </div>
  );
}

function Modal({ open, title, subtitle, children, onClose }) {
  if (!open) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/70 p-2 sm:items-center sm:p-4">
      <div className="absolute inset-0" onClick={onClose} />
      <GlassCard className="relative z-10 max-h-[92vh] w-full max-w-3xl overflow-auto rounded-t-[2rem] p-4 sm:rounded-3xl sm:p-8">
        <div className="mb-6 flex items-start justify-between gap-4">
          <div>
            <h2 className="text-xl font-bold text-white sm:text-2xl">{title}</h2>
            <p className="mt-2 text-sm text-slate-300 sm:text-base">{subtitle}</p>
          </div>
          <button onClick={onClose} className="rounded-2xl border border-white/10 bg-white/5 p-2 text-white transition hover:bg-white/10">
            <X className="h-5 w-5" />
          </button>
        </div>
        {children}
      </GlassCard>
    </div>
  );
}

function PurchaseFlow({ open, plan, onClose, currentUser, onRequireAuth, onComplete }) {
  const [step, setStep] = useState(1);
  const [selectedAddOns, setSelectedAddOns] = useState([]);
  const [setup, setSetup] = useState({ name: "", version: "1.21.2", software: "Paper" });
  const [billing, setBilling] = useState({ fullName: currentUser?.name || "", email: currentUser?.email || "", country: "Austria" });

  useEffect(() => {
    if (open) {
      setBilling({ fullName: currentUser?.name || "", email: currentUser?.email || "", country: "Austria" });
    }
  }, [open, currentUser]);

  if (!open || !plan) return null;

  const toggleAddOn = (key) => {
    setSelectedAddOns((prev) => (prev.includes(key) ? prev.filter((x) => x !== key) : [...prev, key]));
  };

  const selected = addOns.filter((a) => selectedAddOns.includes(a.key));
  const total = plan.price + selected.reduce((sum, item) => sum + item.price, 0);

  const next = () => {
    if (!currentUser) {
      onRequireAuth();
      return;
    }
    if (step === 2 && (!billing.fullName || !billing.email)) return;
    if (step === 3 && !setup.name) return;
    if (step < 3) {
      setStep((s) => s + 1);
      return;
    }
    onComplete({ plan, addons: selected, setup, billing });
    setStep(1);
    setSelectedAddOns([]);
    setSetup({ name: "", version: "1.21.2", software: "Paper" });
  };

  return (
    <Modal open={open} onClose={onClose} title={`Buy ${plan.name}`} subtitle="Choose add-ons, review billing, then create your server through the backend API.">
      <div className="mb-6 flex flex-wrap gap-3">
        {[1, 2, 3].map((n) => (
          <div key={n} className={cn("rounded-full px-4 py-2 text-sm", step === n ? "bg-cyan-400/15 text-cyan-300" : "bg-white/5 text-slate-400")}>
            {n === 1 ? "Extras" : n === 2 ? "Billing" : "Server Setup"}
          </div>
        ))}
      </div>

      {step === 1 && (
        <div className="space-y-6">
          <GlassCard className="p-5">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <div className="text-xl font-bold text-white">{plan.name}</div>
                <div className="mt-2 grid gap-1 text-sm text-slate-300 sm:grid-cols-2">
                  <div>{plan.ram} GB RAM</div>
                  <div>{plan.storage} GB SSD</div>
                  <div>{plan.cores} Cores</div>
                  <div>{plan.cpuPower}% CPU Power</div>
                  <div>{plan.backups} Backups</div>
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm text-slate-400">Base price</div>
                <div className="text-2xl font-bold text-white">{money(plan.price)}</div>
              </div>
            </div>
          </GlassCard>
          <div className="grid gap-4 sm:grid-cols-3">
            {addOns.map((addon) => {
              const active = selectedAddOns.includes(addon.key);
              return (
                <button key={addon.key} onClick={() => toggleAddOn(addon.key)} className={cn("rounded-3xl border p-5 text-left", popClass(), active ? "border-cyan-300/30 bg-cyan-400/10" : "border-white/10 bg-white/5")}>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="font-semibold text-white">{addon.label}</div>
                      <div className="mt-1 text-sm text-slate-400">{money(addon.price)} / month</div>
                    </div>
                    <div className={cn("flex h-7 w-7 items-center justify-center rounded-full border", active ? "border-cyan-300 bg-cyan-300/15 text-cyan-300" : "border-white/10 text-slate-500")}>
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
            <input value={billing.fullName} onChange={(e) => setBilling({ ...billing, fullName: e.target.value })} placeholder="Full name" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
            <input value={billing.email} onChange={(e) => setBilling({ ...billing, email: e.target.value })} placeholder="Email" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
            <input value={billing.country} onChange={(e) => setBilling({ ...billing, country: e.target.value })} placeholder="Country" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
            <div className="rounded-2xl border border-cyan-300/20 bg-cyan-400/10 px-4 py-3 text-sm text-cyan-100">Billing is still UI-only here. The actual backend flow starts when the server is created on the final step.</div>
          </div>
          <GlassCard className="p-5">
            <div className="mb-4 flex items-center gap-2 text-white"><CreditCard className="h-5 w-5 text-cyan-300" /> Order Summary</div>
            <div className="space-y-3 text-sm text-slate-300">
              <div className="flex items-center justify-between"><span>{plan.name}</span><span>{money(plan.price)}</span></div>
              {selected.map((addon) => <div key={addon.key} className="flex items-center justify-between"><span>{addon.label}</span><span>{money(addon.price)}</span></div>)}
              <div className="flex items-center justify-between border-t border-white/10 pt-3 text-base font-semibold text-white"><span>Total</span><span>{money(total)}</span></div>
            </div>
          </GlassCard>
        </div>
      )}

      {step === 3 && (
        <div className="grid gap-4">
          <input value={setup.name} onChange={(e) => setSetup({ ...setup, name: e.target.value })} placeholder="Server name" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
          <select value={setup.software} onChange={(e) => setSetup({ ...setup, software: e.target.value })} className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none">{softwareOptions.map((s) => <option key={s}>{s}</option>)}</select>
          <select value={setup.version} onChange={(e) => setSetup({ ...setup, version: e.target.value })} className="w-full rounded-2xl border border-white/10 bg-slate-900 px-4 py-3 text-white outline-none">{versions.map((v) => <option key={v}>{v}</option>)}</select>
          <div className="rounded-2xl border border-dashed border-cyan-300/25 bg-cyan-400/5 px-4 py-3 text-slate-300">Backend call: bundle/create → server/create. The selected software is mapped to the bundle string.</div>
        </div>
      )}

      <div className="mt-8 flex flex-wrap items-center justify-between gap-4">
        <div className="text-sm text-slate-400">{currentUser ? `Signed in as ${currentUser.email}` : "Create an account or sign in to continue."}</div>
        <div className="flex gap-3">
          {step > 1 && <button onClick={() => setStep((s) => s - 1)} className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 font-semibold text-white">Back</button>}
          <button onClick={next} className="rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950">{!currentUser ? "Sign in to continue" : step === 3 ? "Create Server" : "Continue"}</button>
        </div>
      </div>
    </Modal>
  );
}

function AuthScreen({ mode, onSubmit, setScreen, error, busy }) {
  const signup = mode === "signup";
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  return (
    <div className="mx-auto flex min-h-[calc(100vh-120px)] max-w-7xl items-center px-4 py-8 sm:px-6 sm:py-12 lg:px-16">
      <div className="grid w-full gap-6 lg:grid-cols-[0.95fr_1.05fr] lg:gap-8">
        <GlassCard className="p-5 sm:p-8">
          <div className="mb-5 flex h-14 w-14 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">{signup ? <UserPlus className="h-7 w-7" /> : <LogIn className="h-7 w-7" />}</div>
          <h1 className="text-2xl font-bold text-white sm:text-3xl">{signup ? "Create your account" : "Sign in"}</h1>
          <p className="mt-3 text-slate-300">{signup ? "Create a normal account with email and password. After purchase, only your own servers will appear in your panel." : "Sign in with your email and password to manage your own servers, settings, and purchases."}</p>
          <div className="mt-8 space-y-4">
            {signup && <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Full name / username" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />}
            <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder={signup ? "Email" : "Username or email"} className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
            <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="Password" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
            {error && <div className="rounded-2xl border border-rose-400/20 bg-rose-400/10 px-4 py-3 text-sm text-rose-200">{error}</div>}
            <button disabled={busy} onClick={() => onSubmit({ name, email, password })} className={cn("w-full rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>{busy ? "Please wait..." : signup ? "Create Account" : "Sign In"}</button>
          </div>
          <div className="mt-6 text-sm text-slate-400">{signup ? "Already have an account?" : "Need an account?"} <button onClick={() => setScreen(signup ? "signin" : "signup")} className="font-semibold text-cyan-300">{signup ? "Sign in" : "Create one"}</button></div>
        </GlassCard>
      </div>
    </div>
  );
}

function Sidebar({ items, active, setActive, title, subtitle }) {
  return (
    <GlassCard className="p-4 lg:sticky lg:top-6 lg:h-fit">
      <div className="mb-5 px-2">
        <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">{subtitle}</div>
        <div className="mt-1 text-xl font-bold text-white">{title}</div>
      </div>
      <div className="space-y-2">
        {items.map((item) => {
          const Icon = item.icon;
          return (
            <button key={item.key} onClick={() => setActive(item.key)} className={cn("flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-left", popClass(), active === item.key ? "bg-cyan-400/12 text-white" : "text-slate-300 hover:bg-white/5")}>
              <Icon className="h-5 w-5" />
              <span>{item.label}</span>
            </button>
          );
        })}
      </div>
    </GlassCard>
  );
}

function CustomerDashboard({ currentUser, token, servers, setServers, logout, notices, setNotices }) {
  const [active, setActive] = useState("servers");
  const [selectedServerId, setSelectedServerId] = useState(null);
  const [search, setSearch] = useState("");
  const [modrinthPlugins, setModrinthPlugins] = useState([]);
  const [pluginsLoading, setPluginsLoading] = useState(false);
  const [pluginsError, setPluginsError] = useState("");
  const [backupMap, setBackupMap] = useState({});
  const [statsMap, setStatsMap] = useState({});
  const [serverBusy, setServerBusy] = useState({});

  const ownedServers = useMemo(() => servers.filter((s) => s.ownerId === currentUser.ownerKey), [servers, currentUser.ownerKey]);
  const selectedServer = ownedServers.find((s) => s.server_id === selectedServerId) || ownedServers[0] || null;

  const setBusy = (id, value) => setServerBusy((prev) => ({ ...prev, [id]: value }));
  const pushNotice = (type, text) => setNotices((prev) => [{ id: Date.now() + Math.random(), type, text }, ...prev].slice(0, 5));

  useEffect(() => {
    let ignore = false;
    const loadPlugins = async () => {
      setPluginsLoading(true);
      setPluginsError("");
      try {
        const query = search.trim();
        const facets = encodeURIComponent(JSON.stringify([["project_type:plugin"]]));
        const params = new URLSearchParams({ facets, limit: "12" });
        if (query) params.set("query", query);
        const response = await fetch(`https://api.modrinth.com/v2/search?${params.toString()}`);
        if (!response.ok) throw new Error(`Modrinth request failed with ${response.status}`);
        const data = await response.json();
        const hits = Array.isArray(data?.hits) ? data.hits : [];
        const mapped = hits.map((item) => ({
          id: item.project_id || item.slug || item.title,
          name: item.title,
          source: "Modrinth",
          type: Array.isArray(item.categories) && item.categories.length > 0 ? item.categories.slice(0, 2).join(" · ") : "Plugin",
          description: item.description || "No description available.",
          downloads: item.downloads,
        }));
        if (!ignore) setModrinthPlugins(mapped);
      } catch {
        if (!ignore) {
          setPluginsError("Could not load live Modrinth plugins right now.");
          setModrinthPlugins([]);
        }
      } finally {
        if (!ignore) setPluginsLoading(false);
      }
    };
    if (active === "catalog") loadPlugins();
    return () => { ignore = true; };
  }, [active, search]);

  const runServerAction = async (serverId, action, path) => {
    setBusy(serverId, true);
    try {
      await apiFetch(path, { method: "POST", body: { server_id: serverId }, token });
      setServers((prev) => prev.map((server) => server.server_id === serverId ? { ...server, status: action } : server));
      pushNotice("success", `${action[0].toUpperCase()}${action.slice(1)} request sent for ${serverId}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(serverId, false);
    }
  };

  const fetchBackups = async (serverId) => {
    setBusy(serverId, true);
    try {
      const data = await apiFetch("/server/backup/get", { method: "POST", body: { server_id: serverId }, token });
      setBackupMap((prev) => ({ ...prev, [serverId]: data }));
      pushNotice("success", `Loaded backups for ${serverId}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(serverId, false);
    }
  };

  const createBackup = async (server) => {
    setBusy(server.server_id, true);
    try {
      await apiFetch("/server/backup/create", { method: "POST", body: { server_id: server.server_id, bundle: server.bundle }, token });
      pushNotice("success", `Backup request sent for ${server.name}.`);
      await fetchBackups(server.server_id);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(server.server_id, false);
    }
  };

  const fetchStats = async (serverId) => {
    setBusy(serverId, true);
    try {
      const data = await apiFetch("/server/getStats", { method: "POST", body: { server_id: serverId, key: "overview" }, token });
      setStatsMap((prev) => ({ ...prev, [serverId]: data }));
      pushNotice("success", `Loaded stats for ${serverId}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(serverId, false);
    }
  };

  const deleteServer = async (serverId) => {
    setBusy(serverId, true);
    try {
      await apiFetch("/server/delete", { method: "POST", body: { server_id: serverId }, token });
      setServers((prev) => prev.filter((server) => server.server_id !== serverId));
      pushNotice("success", `Deleted server ${serverId}.`);
    } catch (error) {
      pushNotice("error", error.message);
    } finally {
      setBusy(serverId, false);
    }
  };

  const menu = [
    { key: "servers", label: "My Servers", icon: Server },
    { key: "catalog", label: "Plugins & Mods", icon: Boxes },
    { key: "account", label: "Account", icon: Settings },
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-10 lg:px-16">
      <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">Customer Dashboard</div>
          <h1 className="mt-2 text-3xl font-bold text-white">Welcome back, {currentUser.name}</h1>
          <p className="mt-2 text-slate-300">Only your own servers are visible in this dashboard.</p>
        </div>
        <button onClick={logout} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white", popClass())}>Log out</button>
      </div>

      <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <Sidebar items={menu} active={active} setActive={setActive} title="easy2host Panel" subtitle="User API area" />
        <div className="space-y-6">
          {notices.length > 0 && <GlassCard className="p-4">{notices.map((notice) => <div key={notice.id} className={cn("mb-2 rounded-2xl px-4 py-3 text-sm", notice.type === "error" ? "bg-rose-400/10 text-rose-100" : "bg-cyan-400/10 text-cyan-100")}>{notice.text}</div>)}</GlassCard>}

          {active === "servers" && (
            <>
              <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Owned servers</div><div className="mt-2 text-3xl font-bold text-white">{ownedServers.length}</div></GlassCard>
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Known backups</div><div className="mt-2 text-3xl font-bold text-white">{Object.values(backupMap).length}</div></GlassCard>
                <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">API-ready</div><div className="mt-2 text-lg font-bold text-cyan-300">Compatible</div></GlassCard>
              </div>

              <div className="grid gap-6 2xl:grid-cols-[0.95fr_1.05fr]">
                <GlassCard className="p-5">
                  <h2 className="mb-4 text-xl font-semibold text-white">My Servers</h2>
                  <div className="space-y-3">
                    {ownedServers.length === 0 && <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-6 text-slate-400">You do not have any locally tracked servers yet. Create one through checkout to connect to /bundle/create and /server/create.</div>}
                    {ownedServers.map((server) => (
                      <button key={server.server_id} onClick={() => setSelectedServerId(server.server_id)} className={cn("w-full rounded-2xl border p-4 text-left", popClass(), selectedServer?.server_id === server.server_id ? "border-cyan-300/30 bg-cyan-400/10" : "border-white/10 bg-slate-950/40 hover:bg-white/5")}>
                        <div className="flex items-center justify-between gap-4">
                          <div>
                            <div className="font-semibold text-white">{server.name}</div>
                            <div className="mt-1 text-sm text-slate-400">{server.bundle} · {server.version} · {server.server_id}</div>
                          </div>
                          <div className={cn("rounded-full px-3 py-1 text-xs", server.status === "online" ? "bg-emerald-400/10 text-emerald-300" : "bg-slate-700/60 text-slate-300")}>{server.status || "created"}</div>
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
                          <h2 className="text-2xl font-bold text-white">{selectedServer.name}</h2>
                          <p className="mt-1 text-slate-400">{selectedServer.bundle} · {selectedServer.version} · {selectedServer.server_id}</p>
                        </div>
                        <div className={cn("rounded-full px-3 py-1 text-sm", selectedServer.status === "online" ? "bg-emerald-400/10 text-emerald-300" : "bg-slate-700/60 text-slate-300")}>{selectedServer.status || "created"}</div>
                      </div>

                      <div className="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer.server_id, "online", "/server/start")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-emerald-400/15 px-4 py-3 text-emerald-300 disabled:opacity-60", popClass())}><Play className="h-4 w-4" /> Start</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => runServerAction(selectedServer.server_id, "offline", "/server/stop")} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-400/15 px-4 py-3 text-rose-300 disabled:opacity-60", popClass())}><Square className="h-4 w-4" /> Stop</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => fetchStats(selectedServer.server_id)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-sky-400/15 px-4 py-3 text-sky-300 disabled:opacity-60", popClass())}><Activity className="h-4 w-4" /> Stats</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => createBackup(selectedServer)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-violet-400/15 px-4 py-3 text-violet-300 disabled:opacity-60", popClass())}><HardDrive className="h-4 w-4" /> Backup</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => fetchBackups(selectedServer.server_id)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-white/8 px-4 py-3 text-white disabled:opacity-60", popClass())}><Database className="h-4 w-4" /> Get Backups</button>
                        <button disabled={serverBusy[selectedServer.server_id]} onClick={() => deleteServer(selectedServer.server_id)} className={cn("flex items-center justify-center gap-2 rounded-2xl bg-rose-500/20 px-4 py-3 text-rose-200 disabled:opacity-60", popClass())}><Trash2 className="h-4 w-4" /> Delete</button>
                      </div>

                      <div className="mt-6 grid gap-4 xl:grid-cols-2">
                        <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex items-center gap-2 text-white"><Database className="h-4 w-4 text-cyan-300" /> Stats Response</div>
                          <pre className="max-h-64 overflow-auto text-xs text-slate-300">{JSON.stringify(statsMap[selectedServer.server_id] || { message: "Load stats to see backend response." }, null, 2)}</pre>
                        </div>
                        <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                          <div className="mb-4 flex items-center gap-2 text-white"><HardDrive className="h-4 w-4 text-cyan-300" /> Backups Response</div>
                          <pre className="max-h-64 overflow-auto text-xs text-slate-300">{JSON.stringify(backupMap[selectedServer.server_id] || { message: "Load backups to see backend response." }, null, 2)}</pre>
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

          {active === "catalog" && (
            <GlassCard className="p-6">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Catalog</div>
                  <h2 className="mt-2 text-2xl font-bold text-white">Plugins & Mods</h2>
                  <p className="mt-2 text-slate-400">Search live plugins from Modrinth and browse mods in a clean, simple catalog.</p>
                </div>
                <div className="relative w-full max-w-sm"><Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" /><input value={search} onChange={(e) => setSearch(e.target.value)} className="w-full rounded-2xl border border-white/10 bg-white/5 py-3 pl-11 pr-4 text-white outline-none" placeholder="Search Modrinth plugins or mods" /></div>
              </div>
              <div className="mt-8 grid gap-6 2xl:grid-cols-2">
                <div className="rounded-3xl border border-white/10 bg-slate-950/30 p-5">
                  <div className="mb-4 flex items-center justify-between gap-3"><h3 className="text-xl font-semibold text-white">Plugins</h3><div className="rounded-full bg-cyan-400/10 px-3 py-1 text-xs uppercase tracking-[0.2em] text-cyan-300">Modrinth API</div></div>
                  {pluginsLoading && <div className="rounded-2xl border border-white/10 bg-white/5 p-4 text-sm text-slate-300">Loading live plugins from Modrinth…</div>}
                  {!pluginsLoading && pluginsError && <div className="rounded-2xl border border-amber-400/20 bg-amber-400/10 p-4 text-sm text-amber-100">{pluginsError}</div>}
                  <div className="mt-4 space-y-3">
                    {modrinthPlugins.map((item) => (
                      <div key={item.id} className={cn("rounded-2xl border border-white/10 bg-white/5 p-4", popClass())}>
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <div className="font-semibold text-white">{item.name}</div>
                            <div className="mt-1 text-sm text-slate-400">{item.source} · {item.type}</div>
                            <p className="mt-2 text-sm leading-6 text-slate-300">{item.description}</p>
                            {item.downloads !== undefined && <div className="mt-2 text-xs text-slate-500">{item.downloads.toLocaleString()} downloads</div>}
                          </div>
                          <button className="rounded-xl bg-cyan-400/10 px-3 py-2 text-sm text-cyan-300">Install</button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
                <div className="rounded-3xl border border-white/10 bg-slate-950/30 p-5">
                  <div className="mb-4 flex items-center justify-between gap-3"><h3 className="text-xl font-semibold text-white">Mods</h3><div className="rounded-full bg-white/5 px-3 py-1 text-xs uppercase tracking-[0.2em] text-slate-400">Mocked</div></div>
                  <div className="space-y-3">
                    {fallbackMods.filter((item) => `${item.name} ${item.source} ${item.type}`.toLowerCase().includes(search.toLowerCase())).map((item) => (
                      <div key={item.id} className={cn("flex items-center justify-between gap-4 rounded-2xl border border-white/10 bg-white/5 p-4", popClass())}>
                        <div>
                          <div className="font-semibold text-white">{item.name}</div>
                          <div className="mt-1 text-sm text-slate-400">{item.source} · {item.type}</div>
                          <p className="mt-2 text-sm leading-6 text-slate-300">{item.description}</p>
                        </div>
                        <button className="rounded-xl bg-cyan-400/10 px-3 py-2 text-sm text-cyan-300">Install</button>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </GlassCard>
          )}

          {active === "account" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Account</div>
              <h2 className="mt-2 text-2xl font-bold text-white">Profile & Token</h2>
              <div className="mt-8 grid gap-4 md:grid-cols-3">
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Name</div><div className="mt-2 font-semibold text-white">{currentUser.name}</div></div>
                <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}><div className="text-sm text-slate-400">Email</div><div className="mt-2 font-semibold text-white">{currentUser.email}</div></div>
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

function AdminDashboard({ currentUser, token, notices, setNotices, logout }) {
  const [active, setActive] = useState("overview");
  const [networkIp, setNetworkIp] = useState("");
  const [networkResult, setNetworkResult] = useState(null);
  const [hostMetadata, setHostMetadata] = useState(null);
  const [hostCreateForm, setHostCreateForm] = useState({ ram: "", cores: "" });
  const [hostDeleteId, setHostDeleteId] = useState("");
  const [hostAddForm, setHostAddForm] = useState({ host_server_id: "" });
  const [busy, setBusy] = useState(false);

  const pushNotice = (type, text) => setNotices((prev) => [{ id: Date.now() + Math.random(), type, text }, ...prev].slice(0, 5));

  const adminMenu = [
    { key: "overview", label: "Overview", icon: Crown },
    { key: "network", label: "Network", icon: Server },
    { key: "metadata", label: "Host Metadata", icon: Database },
    { key: "mapping", label: "Add Server To Host", icon: Plus },
  ];

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

  return (
    <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-10 lg:px-16">
      <div className="mb-8 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <div className="text-sm uppercase tracking-[0.25em] text-cyan-300/80">/admin</div>
          <h1 className="mt-2 text-3xl font-bold text-white">Admin Dashboard</h1>
          <p className="mt-2 text-slate-300">Manage infrastructure actions and keep the platform ready for deeper backend features later.</p>
        </div>
        <button onClick={logout} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white", popClass())}>Log out</button>
      </div>

      <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <Sidebar items={adminMenu} active={active} setActive={setActive} title="Admin API Area" subtitle="Documented admin actions" />
        <div className="space-y-6">
          {notices.length > 0 && <GlassCard className="p-4">{notices.map((notice) => <div key={notice.id} className={cn("mb-2 rounded-2xl px-4 py-3 text-sm", notice.type === "error" ? "bg-rose-400/10 text-rose-100" : "bg-cyan-400/10 text-cyan-100")}>{notice.text}</div>)}</GlassCard>}

          {active === "overview" && (
            <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">API Base</div><div className="mt-2 text-lg font-bold text-cyan-300">{API_BASE}</div></GlassCard>
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">Role</div><div className="mt-2 text-3xl font-bold text-white">Admin</div></GlassCard>
              <GlassCard className={cn("p-5", popClass())}><div className="text-sm text-slate-400">JWT present</div><div className="mt-2 text-3xl font-bold text-white">{token ? "Yes" : "No"}</div></GlassCard>
            </div>
          )}

          {active === "network" && (
            <GlassCard className="p-6">
              <div className="text-sm uppercase tracking-[0.2em] text-cyan-300/80">Admin</div>
              <h2 className="mt-2 text-2xl font-bold text-white">Create Network Host Entry</h2>
              <p className="mt-3 text-slate-300">Calls POST /network/create with an IP address.</p>
              <div className="mt-6 flex flex-col gap-4 sm:flex-row">
                <input value={networkIp} onChange={(e) => setNetworkIp(e.target.value)} placeholder="Host IP" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                <button disabled={busy} onClick={async () => {
                  const result = await callAdmin("Network create", () => apiFetch("/network/create", { method: "POST", body: { ip: networkIp }, token }));
                  if (result) setNetworkResult(result);
                }} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Create</button>
              </div>
              <div className="mt-6 rounded-2xl border border-white/10 bg-slate-950/40 p-4"><pre className="overflow-auto text-xs text-slate-300">{JSON.stringify(networkResult || { message: "No result yet." }, null, 2)}</pre></div>
            </GlassCard>
          )}

          {active === "metadata" && (
            <div className="grid gap-6 xl:grid-cols-2">
              <GlassCard className="p-6">
                <h2 className="text-2xl font-bold text-white">Host Metadata</h2>
                <p className="mt-3 text-slate-300">Create, fetch, or delete host metadata using the documented admin endpoints.</p>
                <div className="mt-6 grid gap-4 md:grid-cols-2">
                  <input value={hostCreateForm.ram} onChange={(e) => setHostCreateForm({ ...hostCreateForm, ram: e.target.value })} placeholder="RAM (number only)" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostCreateForm.cores} onChange={(e) => setHostCreateForm({ ...hostCreateForm, cores: e.target.value })} placeholder="CPU cores" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none" />
                  <input value={hostDeleteId} onChange={(e) => setHostDeleteId(e.target.value)} placeholder="Host server ID to delete" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none md:col-span-2" />
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
                <pre className="overflow-auto text-xs text-slate-300">{JSON.stringify(hostMetadata || { message: "No metadata loaded." }, null, 2)}</pre>
              </GlassCard>
            </div>
          )}

          {active === "mapping" && (
            <GlassCard className="p-6">
              <h2 className="text-2xl font-bold text-white">Add Server To Host</h2>
              <p className="mt-3 text-slate-300">Calls POST /host-metadata/add to attach a new server to a host. This endpoint only requires the host server ID in the updated API doc.</p>
              <div className="mt-6 grid gap-4 md:grid-cols-2">
                <input value={hostAddForm.host_server_id} onChange={(e) => setHostAddForm({ ...hostAddForm, host_server_id: e.target.value })} placeholder="Host server ID" className="w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-white outline-none md:col-span-2" />
              </div>
              <div className="mt-6 flex justify-end">
                <button disabled={busy} onClick={() => callAdmin("Host mapping", () => apiFetch("/host-metadata/add", { method: "POST", body: { host_server_id: hostAddForm.host_server_id }, token }))} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-5 py-3 font-semibold text-slate-950 disabled:opacity-60", popClass())}>Attach Server</button>
              </div>
            </GlassCard>
          )}
        </div>
      </div>
    </div>
  );
}

function LandingPage({ setScreen, startPurchase }) {
  return (
    <>
      <section className="relative overflow-hidden px-4 pb-16 pt-8 sm:px-6 sm:pb-20 sm:pt-10 lg:px-16 lg:pb-28 lg:pt-16">
        <div className="absolute inset-0 -z-10"><div className="absolute left-1/2 top-[-20%] h-[28rem] w-[28rem] -translate-x-1/2 rounded-full bg-cyan-500/15 blur-3xl" /><div className="absolute right-[10%] top-[15%] h-[18rem] w-[18rem] rounded-full bg-blue-600/20 blur-3xl" /><div className="absolute left-[5%] top-[35%] h-[16rem] w-[16rem] rounded-full bg-sky-400/10 blur-3xl" /></div>
        <div className="mx-auto grid max-w-7xl items-center gap-8 lg:grid-cols-[1.1fr_0.9fr] lg:gap-10">
          <div>
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-cyan-400/20 bg-cyan-400/10 px-4 py-2 text-sm text-cyan-200">Fast, clean Minecraft hosting</div>
            <div className="mb-6"><LogoMark full /></div>
            <h1 className="max-w-3xl text-4xl font-black tracking-tight text-white sm:text-5xl lg:text-7xl">Start Hosting Today <span className="bg-gradient-to-r from-cyan-300 via-sky-300 to-blue-400 bg-clip-text text-transparent">with Easy2host.</span></h1>
            <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300 sm:mt-6 sm:text-lg sm:leading-8">Create your account, choose the plan that fits your server, add optional upgrades, complete checkout, and finish with your server setup. After purchase, your server is deployed into your personal dashboard, where only you can manage it.</p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:gap-4"><button onClick={() => setScreen("signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-6 py-3 font-semibold text-slate-950", popClass())}>Create Account</button><a href="#plans" className={cn("rounded-2xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-white", popClass())}>Explore Plans</a></div>
            <div className="mt-10 grid max-w-xl grid-cols-1 gap-4 text-sm text-slate-300 sm:grid-cols-3">
              <GlassCard className={cn("p-4", popClass())}><div className="text-2xl font-bold text-white">1.8 → 1.21.2</div><div className="mt-1">Old and new versions supported</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="text-2xl font-bold text-white">Extras</div><div className="mt-1">Add RAM and SSD during checkout</div></GlassCard>
              <GlassCard className={cn("p-4", popClass())}><div className="text-2xl font-bold text-white">Private</div><div className="mt-1">Only your own servers show up</div></GlassCard>
            </div>
          </div>
          <GlassCard className={cn("overflow-hidden p-6 lg:p-7", popClass())}>
            <div className="mb-4 flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-400">How it works</p>
                <h3 className="text-xl font-semibold text-white">Purchase flow preview</h3>
              </div>
              <div className="rounded-xl border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-sm text-cyan-200">3 steps</div>
            </div>
            <div className="space-y-4">
              <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}>
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">1</div>
                  <div>
                    <div className="text-lg font-semibold text-white">Choose your hosting plan</div>
                    <div className="mt-1 text-sm text-slate-400">Select Starter, Standard, or Premium, then add extra RAM or SSD if you want more resources.</div>
                  </div>
                </div>
              </div>
              <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}>
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">2</div>
                  <div>
                    <div className="text-lg font-semibold text-white">Complete checkout</div>
                    <div className="mt-1 text-sm text-slate-400">Review your order and enter billing details in a dedicated checkout step before the server is created.</div>
                  </div>
                </div>
              </div>
              <div className={cn("rounded-2xl border border-white/10 bg-slate-950/40 p-4", popClass())}>
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">3</div>
                  <div>
                    <div className="text-lg font-semibold text-white">Configure your server</div>
                    <div className="mt-1 text-sm text-slate-400">Choose the server name, software type, and Minecraft version, then manage it from your own dashboard.</div>
                  </div>
                </div>
              </div>
            </div>
          </GlassCard>
        </div>
      </section>

      <section id="plans" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <div className="max-w-2xl"><p className="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-cyan-300/80">Plans</p><h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">Choose a plan, then create a server</h2><p className="mt-4 text-base leading-7 text-slate-300">Choose a plan, customize it with add-ons, and finish with your server setup before deployment.</p></div>
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {plans.map((plan) => (
            <button key={plan.name} onClick={() => startPurchase(plan)} className={cn("relative rounded-3xl border p-6 text-left", popClass(), plan.featured ? "border-cyan-300/30 bg-cyan-400/10 shadow-cyan-500/10" : "border-white/10 bg-white/5")}>
              {plan.featured && <div className="absolute right-5 top-5 rounded-full bg-gradient-to-r from-cyan-300 to-blue-500 px-3 py-1 text-xs font-bold uppercase tracking-[0.2em] text-slate-950">Most popular</div>}
              <h3 className="text-2xl font-bold text-white">{plan.name}</h3>
              <div className="mt-6 space-y-2 text-slate-300">
                <div>{plan.ram} GB RAM</div>
                <div>{plan.storage} GB SSD</div>
                <div>{plan.cores} Cores</div>
                <div>{plan.cpuPower}% CPU Power</div>
                <div>{plan.backups} Backups</div>
              </div>
              <div className="mt-8 text-3xl font-black text-white">{money(plan.price)} <span className="text-base font-medium text-slate-400">/ month</span></div>
              <div className="mt-8 inline-flex items-center gap-2 rounded-2xl bg-white/8 px-4 py-3 font-semibold text-white">Select plan <Plus className="h-4 w-4" /></div>
            </button>
          ))}
        </div>

        <div className="mt-12">
          <div className="mb-5 flex items-center justify-between gap-4">
            <div>
              <h3 className="text-2xl font-bold text-white">Available upgrades</h3>
              <p className="mt-2 text-slate-300">These extra upgrade windows are shown again so users can clearly see what can be added before checkout.</p>
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
                  <div>Can be selected during checkout</div>
                </div>
                <button onClick={() => startPurchase(plans[1])} className={cn("mt-6 w-full rounded-2xl border border-white/10 bg-white/5 px-4 py-3 font-semibold text-white", popClass())}>Choose with a plan</button>
              </GlassCard>
            ))}
          </div>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <div className="max-w-2xl">
          <p className="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-cyan-300/80">Features</p>
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">Everything you need to manage your server cleanly</h2>
          <p className="mt-4 text-base leading-7 text-slate-300">easy2host is built to feel simple for beginners while still giving you the controls that matter.</p>
        </div>
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {[
            [LayoutDashboard, "Clean dashboard", "A clear panel for server status, actions, backups, and account management."],
            [Server, "Easy server control", "Start, stop, back up, and manage your server from one simple interface."],
            [Boxes, "Plugins and mods", "Browse live plugins from Modrinth and keep your setup flexible."],
            [Shield, "Account-based access", "Each account only sees and manages its own servers and resources."],
            [Settings, "Version and software choice", "Pick Vanilla, Fabric, Bukkit, or Paper together with your preferred Minecraft version."],
            [CreditCard, "Simple purchase flow", "Choose a plan, add upgrades, review checkout, and finish with server setup."],
          ].map(([Icon, title, text]) => (
            <GlassCard key={title} className={cn("p-6", popClass())}>
              <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-400/10 text-cyan-300">
                <Icon className="h-6 w-6" />
              </div>
              <h3 className="mt-5 text-xl font-semibold text-white">{title}</h3>
              <p className="mt-3 leading-7 text-slate-300">{text}</p>
            </GlassCard>
          ))}
        </div>
      </section>

      <section id="faq" className="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-20 lg:px-16">
        <div className="max-w-2xl">
          <p className="mb-3 text-sm font-semibold uppercase tracking-[0.2em] text-cyan-300/80">FAQ</p>
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">Frequently asked questions</h2>
          <p className="mt-4 text-base leading-7 text-slate-300">Clear answers for the things most people want to know before getting started.</p>
        </div>
        <div className="mt-10 grid gap-4">
          {[
            ["How fast can I create a server?", "Once your account is ready, you can choose a plan, complete checkout, and finish server setup in just a few steps."],
            ["Can I change the Minecraft version later?", "Yes. The platform is designed around editable server settings and flexible version selection."],
            ["Which server software can I choose?", "You can choose between Vanilla, Fabric, Bukkit, and Paper during setup."],
            ["Can I add extra resources?", "Yes. During checkout you can add extra RAM and SSD upgrades before the server is created."],
            ["Will other users be able to see my server?", "No. Your account dashboard only shows the servers connected to your own account."],
            ["Does the panel support plugins and mods?", "Yes. The interface includes a plugin and mod catalog, including live plugin search through Modrinth."],
          ].map(([question, answer]) => (
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
                <h2 className="mt-3 text-3xl font-bold text-white sm:text-4xl">Create your account and launch your server today</h2>
                <p className="mt-4 max-w-2xl text-slate-300">Pick your plan, configure your server, and manage everything from one clean dashboard.</p>
              </div>
              <button onClick={() => setScreen("signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-6 py-3 font-semibold text-slate-950", popClass())}>Create Account</button>
            </div>
          </GlassCard>
        </div>
      </section>
    </>
  );
}

export default function Easy2HostPlatformPrototype() {
  const [screen, setScreen] = useState("landing");
  const [mobileOpen, setMobileOpen] = useState(false);
  const [token, setToken] = useState(() => localStorage.getItem("easy2host_token") || "");
  const [currentUser, setCurrentUser] = useState(() => safeJsonParse(localStorage.getItem("easy2host_user"), null));
  const [servers, setServers] = useState(() => safeJsonParse(localStorage.getItem("easy2host_servers"), []));
  const [authBusy, setAuthBusy] = useState(false);
  const [authError, setAuthError] = useState("");
  const [purchasePlan, setPurchasePlan] = useState(null);
  const [notices, setNotices] = useState([]);

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

  const startPurchase = (plan) => setPurchasePlan(plan);

  const logout = () => {
    setToken("");
    setCurrentUser(null);
    setScreen("landing");
  };

  const buildUserFromToken = (loginName, jwt) => {
    const decoded = decodeJwt(jwt) || {};
    const email = decoded.email || decoded.username || decoded.sub || loginName;
    const normalizedEmail = String(email || loginName || "").toLowerCase();
    const role = ADMIN_EMAILS.includes(normalizedEmail) || decoded.role === "admin" || decoded.type === "admin" ? "admin" : "user";
    return {
      name: decoded.username || decoded.name || loginName,
      email: normalizedEmail,
      role,
      ownerKey: normalizedEmail,
      raw: decoded,
    };
  };

  const handleSignUp = async ({ name, email, password }) => {
    setAuthBusy(true);
    setAuthError("");
    try {
      await apiFetch("/auth/user/register", {
        method: "POST",
        body: {
          username: name || email,
          password,
          email,
          type: "user",
          jwt: "",
        },
      });
      setScreen("signin");
      setNotices((prev) => [{ id: Date.now(), type: "success", text: "Account created. You can sign in now." }, ...prev].slice(0, 5));
    } catch (error) {
      setAuthError(error.message);
    } finally {
      setAuthBusy(false);
    }
  };

  const handleSignIn = async ({ email, password }) => {
    setAuthBusy(true);
    setAuthError("");
    try {
      const data = await apiFetch("/auth/user/login", { method: "POST", body: { username: email, password } });
      const nextToken = data.token || "";
      const user = buildUserFromToken(email, nextToken);
      setToken(nextToken);
      setCurrentUser(user);
      setScreen(user.role === "admin" ? "admin" : "dashboard");
      setNotices((prev) => [{ id: Date.now(), type: "success", text: `Signed in as ${user.email}.` }, ...prev].slice(0, 5));
    } catch (error) {
      setAuthError(error.message);
    } finally {
      setAuthBusy(false);
    }
  };

  const completePurchase = async (order) => {
    if (!currentUser || !token) {
      setScreen("signin");
      return;
    }
    try {
      const bundleName = `${order.setup.software.toLowerCase()}-${order.setup.version}`;
      await apiFetch("/bundle/create", { method: "POST", body: { bundle: bundleName }, token });
      await apiFetch("/bundle/add", { method: "POST", body: { bundle: bundleName }, token });
      const created = await apiFetch("/server/create", {
        method: "POST",
        body: { version: order.setup.version, bundle: bundleName },
        token,
      });
      const extraRam = order.addons.reduce((sum, a) => sum + a.ram, 0);
      const extraStorage = order.addons.reduce((sum, a) => sum + a.storage, 0);
      setServers((prev) => [
        ...prev,
        {
          ownerId: currentUser.ownerKey,
          server_id: created.server_id,
          port: created.port,
          name: order.setup.name,
          version: order.setup.version,
          bundle: bundleName,
          software: order.setup.software,
          status: "created",
          plan: order.plan.name,
          ram: order.plan.ram + extraRam,
          storage: order.plan.storage + extraStorage,
          cores: order.plan.cores,
          cpuPower: order.plan.cpuPower,
          backups: order.plan.backups,
        },
      ]);
      setPurchasePlan(null);
      setScreen("dashboard");
      setNotices((prev) => [{ id: Date.now(), type: "success", text: `Server ${order.setup.name} created with ID ${created.server_id}.` }, ...prev].slice(0, 5));
    } catch (error) {
      setNotices((prev) => [{ id: Date.now(), type: "error", text: error.message }, ...prev].slice(0, 5));
    }
  };

  return (
    <div className="min-h-screen bg-[#07111f] text-white">
      <div className="fixed inset-0 -z-10 bg-[radial-gradient(circle_at_top,rgba(56,189,248,0.08),transparent_28%),linear-gradient(180deg,#08111f_0%,#09121c_45%,#070d17_100%)]" />
      <header className="sticky top-0 z-40 border-b border-white/10 bg-slate-950/70 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-4 sm:px-6 lg:px-16">
          <button onClick={() => setScreen("landing")} className="min-w-0 text-left"><LogoMark full /></button>
          <nav className="hidden items-center gap-8 md:flex">
            <button onClick={() => setScreen("landing")} className="text-sm font-medium text-slate-300 transition hover:text-white">Home</button>
            <a href="#plans" className="text-sm font-medium text-slate-300 transition hover:text-white">Plans</a>
            <a href="#features" className="text-sm font-medium text-slate-300 transition hover:text-white">Features</a>
            <a href="#faq" className="text-sm font-medium text-slate-300 transition hover:text-white">FAQ</a>
            {currentUser && currentUser.role !== "admin" && <button onClick={() => setScreen("dashboard")} className="text-sm font-medium text-slate-300 transition hover:text-white">Dashboard</button>}
            {currentUser?.role === "admin" && <button onClick={() => setScreen("admin")} className="text-sm font-medium text-slate-300 transition hover:text-white">/admin</button>}
          </nav>
          <div className="hidden items-center gap-2 lg:gap-3 md:flex">
            {currentUser ? (
              <>
                <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-slate-300">{currentUser.email}</div>
                <button onClick={() => setScreen(currentUser.role === "admin" ? "admin" : "dashboard")} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 font-medium text-white", popClass())}>{currentUser.role === "admin" ? "Admin Panel" : "Dashboard"}</button>
              </>
            ) : (
              <>
                <button onClick={() => setScreen("signin")} className={cn("rounded-2xl border border-white/10 bg-white/5 px-4 py-2.5 font-medium text-white", popClass())}>Sign In</button>
                <button onClick={() => setScreen("signup")} className={cn("rounded-2xl bg-gradient-to-r from-cyan-300 to-blue-500 px-4 py-2.5 font-semibold text-slate-950", popClass())}>Create Account</button>
              </>
            )}
          </div>
          <button onClick={() => setMobileOpen((v) => !v)} className="rounded-2xl border border-white/10 bg-white/5 p-2.5 text-white md:hidden">{mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}</button>
        </div>
        {mobileOpen && <div className="border-t border-white/10 px-4 py-4 sm:px-6 md:hidden"><div className="flex flex-col gap-3"><button onClick={() => { setScreen("landing"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Home</button><a href="#plans" onClick={() => setMobileOpen(false)} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Plans</a><a href="#features" onClick={() => setMobileOpen(false)} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Features</a><a href="#faq" onClick={() => setMobileOpen(false)} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">FAQ</a>{!currentUser && <><button onClick={() => { setScreen("signin"); setMobileOpen(false); }} className="rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-left text-white">Sign In</button><button onClick={() => { setScreen("signup"); setMobileOpen(false); }} className="rounded-xl bg-gradient-to-r from-cyan-300 to-blue-500 px-3 py-2 text-left font-semibold text-slate-950">Create Account</button></>}{currentUser && currentUser.role !== "admin" && <button onClick={() => { setScreen("dashboard"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">Dashboard</button>}{currentUser?.role === "admin" && <button onClick={() => { setScreen("admin"); setMobileOpen(false); }} className="rounded-xl px-3 py-2 text-left text-slate-300 hover:bg-white/5 hover:text-white">/admin</button>}</div></div>}
      </header>

      {screen === "landing" && <LandingPage setScreen={setScreen} startPurchase={startPurchase} />}
      {screen === "signup" && <AuthScreen mode="signup" onSubmit={handleSignUp} setScreen={setScreen} error={authError} busy={authBusy} />}
      {screen === "signin" && <AuthScreen mode="signin" onSubmit={handleSignIn} setScreen={setScreen} error={authError} busy={authBusy} />}
      {screen === "dashboard" && currentUser && currentUser.role !== "admin" && <CustomerDashboard currentUser={currentUser} token={token} servers={servers} setServers={setServers} logout={logout} notices={notices} setNotices={setNotices} />}
      {screen === "admin" && currentUser?.role === "admin" && <AdminDashboard currentUser={currentUser} token={token} notices={notices} setNotices={setNotices} logout={logout} />}

      <PurchaseFlow open={!!purchasePlan} plan={purchasePlan} onClose={() => setPurchasePlan(null)} currentUser={currentUser} onRequireAuth={() => setScreen("signin")} onComplete={completePurchase} />

      <footer className="border-t border-white/10 px-4 py-8 text-center text-sm text-slate-400 sm:px-6 lg:px-16">easy2host prototype — clean frontend experience with compatibility hooks for auth, server lifecycle, backups, and admin infrastructure actions.</footer>
    </div>
  );
}
