"use client";

import { useState } from "react";

// ── Types ──────────────────────────────────────────────────────────────────
type TimeRange = "daily" | "weekly" | "monthly";
type Condition =
  | "OPTIMAL"
  | "EYE_STRAIN_RISK"
  | "POSTURE_RISK"
  | "DISTRACTED"
  | "AWAY";
type NavItem = "dashboard" | "analytics" | "settings";

interface Session {
  id: string;
  started_at: string;
  ended_at: string;
  duration_sec: number;
  dominant_condition: Condition;
}

interface HourBar {
  hour: number;
  minutes: number;
  condition: Condition;
}

// ── Constants ──────────────────────────────────────────────────────────────
const CONDITION_META: Record<
  Condition,
  { label: string; color: string; bg: string; bar: string }
> = {
  OPTIMAL: {
    label: "Optimal",
    color: "#1D6A3A",
    bg: "#E6F4EC",
    bar: "#22C55E",
  },
  EYE_STRAIN_RISK: {
    label: "Eye Strain Risk",
    color: "#92400E",
    bg: "#FEF3C7",
    bar: "#F59E0B",
  },
  POSTURE_RISK: {
    label: "Posture Risk",
    color: "#92400E",
    bg: "#FEF3C7",
    bar: "#FBBF24",
  },
  DISTRACTED: {
    label: "Distracted",
    color: "#9A3412",
    bg: "#FFF0E6",
    bar: "#F97316",
  },
  AWAY: { label: "Away", color: "#374151", bg: "#F3F4F6", bar: "#D1D5DB" },
};

// ── Mock Data ──────────────────────────────────────────────────────────────
const MOCK_HOURLY: HourBar[] = [
  { hour: 8, minutes: 12, condition: "OPTIMAL" },
  { hour: 9, minutes: 58, condition: "OPTIMAL" },
  { hour: 10, minutes: 45, condition: "EYE_STRAIN_RISK" },
  { hour: 11, minutes: 30, condition: "OPTIMAL" },
  { hour: 12, minutes: 5, condition: "DISTRACTED" },
  { hour: 13, minutes: 0, condition: "AWAY" },
  { hour: 14, minutes: 38, condition: "POSTURE_RISK" },
  { hour: 15, minutes: 55, condition: "OPTIMAL" },
  { hour: 16, minutes: 20, condition: "EYE_STRAIN_RISK" },
  { hour: 17, minutes: 10, condition: "OPTIMAL" },
  { hour: 18, minutes: 0, condition: "AWAY" },
  { hour: 19, minutes: 42, condition: "OPTIMAL" },
  { hour: 20, minutes: 28, condition: "DISTRACTED" },
];

const MOCK_SESSIONS: Session[] = [
  {
    id: "SSN-001",
    started_at: "08:02",
    ended_at: "10:48",
    duration_sec: 9960,
    dominant_condition: "OPTIMAL",
  },
  {
    id: "SSN-002",
    started_at: "11:05",
    ended_at: "12:33",
    duration_sec: 5280,
    dominant_condition: "EYE_STRAIN_RISK",
  },
  {
    id: "SSN-003",
    started_at: "14:10",
    ended_at: "16:05",
    duration_sec: 7500,
    dominant_condition: "POSTURE_RISK",
  },
  {
    id: "SSN-004",
    started_at: "19:20",
    ended_at: "20:52",
    duration_sec: 5520,
    dominant_condition: "OPTIMAL",
  },
];

const MOCK_BREAKDOWN: { condition: Condition; percent: number }[] = [
  { condition: "OPTIMAL", percent: 54 },
  { condition: "EYE_STRAIN_RISK", percent: 18 },
  { condition: "POSTURE_RISK", percent: 11 },
  { condition: "DISTRACTED", percent: 9 },
  { condition: "AWAY", percent: 8 },
];

const PEAK_HOURS = ["09:00 – 10:00", "15:00 – 16:00", "19:00 – 20:00"];

// ── Helpers ────────────────────────────────────────────────────────────────
function fmtDuration(sec: number): string {
  const h = Math.floor(sec / 3600);
  const m = Math.floor((sec % 3600) / 60);
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

function fmtHour(h: number): string {
  return `${String(h).padStart(2, "0")}:00`;
}

// ── Nav helpers ────────────────────────────────────────────────────────────
function NavIcon({ id }: { id: NavItem }) {
  if (id === "dashboard")
    return (
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <rect x="3" y="3" width="7" height="7" rx="1" />
        <rect x="14" y="3" width="7" height="7" rx="1" />
        <rect x="3" y="14" width="7" height="7" rx="1" />
        <rect x="14" y="14" width="7" height="7" rx="1" />
      </svg>
    );
  if (id === "analytics")
    return (
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      >
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
      </svg>
    );
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="3" />
      <path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14" />
    </svg>
  );
}

function Logo() {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 9 }}>
      <svg width="30" height="36" viewBox="0 0 56 68" fill="none">
        <path
          d="M28 4C15.85 4 6 13.85 6 26c0 8.4 4.6 15.7 11.4 19.6V52h21.2v-6.4C45.4 41.7 50 34.4 50 26 50 13.85 40.15 4 28 4z"
          fill="#F5A823"
        />
        <rect x="17.4" y="54" width="21.2" height="8" rx="4" fill="#F5A823" />
        <ellipse cx="28" cy="25" rx="10" ry="6.5" fill="white" />
        <circle cx="28" cy="25" r="4" fill="#1a1a18" />
        <circle cx="29.6" cy="23.4" r="1.2" fill="white" />
        <line
          x1="28"
          y1="14"
          x2="27.2"
          y2="17.2"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
        <line
          x1="22.8"
          y1="15.8"
          x2="23.8"
          y2="18.8"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
        <line
          x1="33.2"
          y1="15.8"
          x2="32.2"
          y2="18.8"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
        <line
          x1="28"
          y1="36"
          x2="27.2"
          y2="32.8"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
        <line
          x1="22.8"
          y1="34.2"
          x2="23.8"
          y2="31.2"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
        <line
          x1="33.2"
          y1="34.2"
          x2="32.2"
          y2="31.2"
          stroke="white"
          strokeWidth="2"
          strokeLinecap="round"
        />
      </svg>
      <span
        style={{
          fontSize: "1rem",
          fontWeight: 700,
          letterSpacing: "-0.03em",
          color: "#1a1a18",
        }}
      >
        StareDesk
      </span>
    </div>
  );
}

// ── TimelineChart ──────────────────────────────────────────────────────────
function TimelineChart({ data, range }: { data: HourBar[]; range: TimeRange }) {
  const maxMin = Math.max(...data.map((d) => d.minutes), 1);

  return (
    <div
      style={{
        background: "#fff",
        border: "1.5px solid #e8e8e2",
        borderRadius: 16,
        padding: "20px 22px",
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: 20,
        }}
      >
        <p
          style={{
            fontSize: "0.75rem",
            fontWeight: 500,
            color: "#888880",
            textTransform: "uppercase",
            letterSpacing: "0.06em",
          }}
        >
          Focus timeline
        </p>
        <span style={{ fontSize: "0.75rem", color: "#aaa" }}>
          {range === "daily"
            ? "Today"
            : range === "weekly"
              ? "This week"
              : "This month"}
        </span>
      </div>

      {/* Bar chart */}
      <div
        style={{ display: "flex", alignItems: "flex-end", gap: 6, height: 120 }}
      >
        {data.map((bar) => {
          const heightPct = (bar.minutes / maxMin) * 100;
          const meta = CONDITION_META[bar.condition];
          return (
            <div
              key={bar.hour}
              title={`${fmtHour(bar.hour)} — ${bar.minutes}m (${meta.label})`}
              style={{
                flex: 1,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                gap: 6,
                height: "100%",
                justifyContent: "flex-end",
                cursor: "default",
              }}
            >
              <div
                style={{
                  width: "100%",
                  height: `${heightPct}%`,
                  background: meta.bar,
                  borderRadius: "4px 4px 0 0",
                  minHeight: bar.minutes > 0 ? 4 : 0,
                  transition: "height 0.3s ease",
                }}
              />
              <span
                style={{
                  fontSize: "0.6rem",
                  color: "#bbb",
                  whiteSpace: "nowrap",
                }}
              >
                {bar.hour}
              </span>
            </div>
          );
        })}
      </div>

      {/* Legend */}
      <div
        style={{
          display: "flex",
          flexWrap: "wrap",
          gap: "6px 14px",
          marginTop: 16,
        }}
      >
        {(
          [
            "OPTIMAL",
            "EYE_STRAIN_RISK",
            "POSTURE_RISK",
            "DISTRACTED",
            "AWAY",
          ] as Condition[]
        ).map((c) => (
          <div
            key={c}
            style={{ display: "flex", alignItems: "center", gap: 5 }}
          >
            <span
              style={{
                width: 8,
                height: 8,
                borderRadius: 2,
                background: CONDITION_META[c].bar,
                display: "inline-block",
                flexShrink: 0,
              }}
            />
            <span style={{ fontSize: "0.7rem", color: "#888880" }}>
              {CONDITION_META[c].label}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

// ── ConditionBreakdown ─────────────────────────────────────────────────────
function ConditionBreakdown({ data }: { data: typeof MOCK_BREAKDOWN }) {
  // Simple segmented horizontal bar
  return (
    <div
      style={{
        background: "#fff",
        border: "1.5px solid #e8e8e2",
        borderRadius: 16,
        padding: "20px 22px",
      }}
    >
      <p
        style={{
          fontSize: "0.75rem",
          fontWeight: 500,
          color: "#888880",
          textTransform: "uppercase",
          letterSpacing: "0.06em",
          marginBottom: 18,
        }}
      >
        Condition breakdown
      </p>

      {/* Segmented bar */}
      <div
        style={{
          display: "flex",
          borderRadius: 6,
          overflow: "hidden",
          height: 12,
          marginBottom: 18,
        }}
      >
        {data.map((d) => (
          <div
            key={d.condition}
            title={`${CONDITION_META[d.condition].label}: ${d.percent}%`}
            style={{
              width: `${d.percent}%`,
              background: CONDITION_META[d.condition].bar,
            }}
          />
        ))}
      </div>

      {/* Rows */}
      <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        {data.map((d) => {
          const meta = CONDITION_META[d.condition];
          return (
            <div
              key={d.condition}
              style={{ display: "flex", alignItems: "center", gap: 10 }}
            >
              <span
                style={{
                  width: 8,
                  height: 8,
                  borderRadius: 2,
                  background: meta.bar,
                  flexShrink: 0,
                }}
              />
              <span style={{ fontSize: "0.82rem", color: "#1a1a18", flex: 1 }}>
                {meta.label}
              </span>
              <div
                style={{
                  width: 80,
                  background: "#f3f3ef",
                  borderRadius: 4,
                  height: 5,
                  overflow: "hidden",
                }}
              >
                <div
                  style={{
                    width: `${d.percent}%`,
                    height: "100%",
                    background: meta.bar,
                    borderRadius: 4,
                  }}
                />
              </div>
              <span
                style={{
                  fontSize: "0.8rem",
                  fontWeight: 600,
                  color: "#1a1a18",
                  minWidth: 32,
                  textAlign: "right",
                }}
              >
                {d.percent}%
              </span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ── PeakHoursCard ──────────────────────────────────────────────────────────
function PeakHoursCard({ hours }: { hours: string[] }) {
  return (
    <div
      style={{
        background: "#fff",
        border: "1.5px solid #e8e8e2",
        borderRadius: 16,
        padding: "20px 22px",
      }}
    >
      <p
        style={{
          fontSize: "0.75rem",
          fontWeight: 500,
          color: "#888880",
          textTransform: "uppercase",
          letterSpacing: "0.06em",
          marginBottom: 14,
        }}
      >
        Peak hours
      </p>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {hours.map((h, i) => (
          <div
            key={h}
            style={{ display: "flex", alignItems: "center", gap: 12 }}
          >
            <span
              style={{
                width: 22,
                height: 22,
                borderRadius: "50%",
                background: i === 0 ? "#F5A823" : "#f3f3ef",
                color: i === 0 ? "#1a1a18" : "#888880",
                fontSize: "0.7rem",
                fontWeight: 700,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                flexShrink: 0,
              }}
            >
              {i + 1}
            </span>
            <span
              style={{
                fontSize: "0.88rem",
                fontWeight: i === 0 ? 600 : 400,
                color: "#1a1a18",
              }}
            >
              {h}
            </span>
            {i === 0 && (
              <span
                style={{
                  marginLeft: "auto",
                  fontSize: "0.68rem",
                  fontWeight: 500,
                  padding: "2px 8px",
                  borderRadius: 20,
                  background: "#FEF3C7",
                  color: "#92400E",
                }}
              >
                Best
              </span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

// ── SessionHistoryTable ────────────────────────────────────────────────────
function SessionHistoryTable({ sessions }: { sessions: Session[] }) {
  return (
    <div
      style={{
        background: "#fff",
        border: "1.5px solid #e8e8e2",
        borderRadius: 16,
        padding: "20px 22px",
      }}
    >
      <p
        style={{
          fontSize: "0.75rem",
          fontWeight: 500,
          color: "#888880",
          textTransform: "uppercase",
          letterSpacing: "0.06em",
          marginBottom: 16,
        }}
      >
        Session history
      </p>
      <div style={{ display: "flex", flexDirection: "column", gap: 0 }}>
        {/* Header */}
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "1fr 80px 90px 120px",
            padding: "0 0 8px",
            borderBottom: "1px solid #e8e8e2",
            gap: 8,
          }}
        >
          {["Time", "Duration", "Condition", ""].map((h) => (
            <span
              key={h}
              style={{
                fontSize: "0.68rem",
                fontWeight: 500,
                color: "#aaa",
                textTransform: "uppercase",
                letterSpacing: "0.05em",
              }}
            >
              {h}
            </span>
          ))}
        </div>
        {sessions.map((s, i) => {
          const meta = CONDITION_META[s.dominant_condition];
          return (
            <div
              key={s.id}
              style={{
                display: "grid",
                gridTemplateColumns: "1fr 80px 90px 120px",
                padding: "12px 0",
                borderBottom:
                  i < sessions.length - 1 ? "1px solid #f0f0ec" : "none",
                alignItems: "center",
                gap: 8,
              }}
            >
              <span
                style={{
                  fontSize: "0.82rem",
                  color: "#1a1a18",
                  fontWeight: 500,
                }}
              >
                {s.started_at} – {s.ended_at}
              </span>
              <span style={{ fontSize: "0.82rem", color: "#555550" }}>
                {fmtDuration(s.duration_sec)}
              </span>
              <span
                style={{
                  fontSize: "0.72rem",
                  fontWeight: 500,
                  padding: "3px 9px",
                  borderRadius: 20,
                  background: meta.bg,
                  color: meta.color,
                  display: "inline-block",
                }}
              >
                {meta.label}
              </span>
              <span style={{ fontSize: "0.72rem", color: "#bbb" }}>{s.id}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ── Main Analytics Page ────────────────────────────────────────────────────
export default function AnalyticsPage() {
  const [activeNav, setActiveNav] = useState<NavItem>("analytics");
  const [range, setRange] = useState<TimeRange>("daily");

  const navItems: { id: NavItem; label: string }[] = [
    { id: "dashboard", label: "Dashboard" },
    { id: "analytics", label: "Analytics" },
    { id: "settings", label: "Settings" },
  ];

  const totalFocusMin = MOCK_HOURLY.reduce((acc, d) => acc + d.minutes, 0);
  const totalFocusLabel = `${Math.floor(totalFocusMin / 60)}h ${totalFocusMin % 60}m`;

  return (
    <div className="root">
      {/* ── Sidebar ── */}
      <aside className="sidebar">
        <div className="sidebar-logo">
          <Logo />
        </div>
        <nav className="sidebar-nav">
          {navItems.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${activeNav === item.id ? "nav-active" : ""}`}
              onClick={() => setActiveNav(item.id)}
            >
              <NavIcon id={item.id} />
              <span>{item.label}</span>
            </button>
          ))}
        </nav>
        <div className="sidebar-footer">
          <button className="logout-btn">
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
              <polyline points="16 17 21 12 16 7" />
              <line x1="21" y1="12" x2="9" y2="12" />
            </svg>
            <span>Sign out</span>
          </button>
        </div>
      </aside>

      {/* ── Main ── */}
      <main className="main">
        <header className="topbar">
          <div className="topbar-logo-mobile">
            <Logo />
          </div>
          <div style={{ flex: 1 }}>
            <h1 className="page-title">Analytics</h1>
            <p className="page-subtitle">Monday, 28 April 2026</p>
          </div>
          {/* Range toggle */}
          <div className="range-toggle">
            {(["daily", "weekly", "monthly"] as TimeRange[]).map((r) => (
              <button
                key={r}
                className={`toggle-btn ${range === r ? "toggle-active" : ""}`}
                onClick={() => setRange(r)}
              >
                {r.charAt(0).toUpperCase() + r.slice(1)}
              </button>
            ))}
          </div>
        </header>

        <div className="content">
          {/* Summary strip */}
          <div className="summary-strip">
            <div className="strip-card">
              <p className="strip-label">Total focus</p>
              <p className="strip-value">{totalFocusLabel}</p>
            </div>
            <div className="strip-card">
              <p className="strip-label">Sessions</p>
              <p className="strip-value">{MOCK_SESSIONS.length}</p>
            </div>
            <div className="strip-card">
              <p className="strip-label">Optimal rate</p>
              <p className="strip-value">
                {MOCK_BREAKDOWN.find((d) => d.condition === "OPTIMAL")?.percent}
                %
              </p>
            </div>
            <div className="strip-card">
              <p className="strip-label">Avg session</p>
              <p className="strip-value">
                {fmtDuration(
                  Math.round(
                    MOCK_SESSIONS.reduce((a, s) => a + s.duration_sec, 0) /
                      MOCK_SESSIONS.length,
                  ),
                )}
              </p>
            </div>
          </div>

          {/* Timeline full-width */}
          <TimelineChart data={MOCK_HOURLY} range={range} />

          {/* Two-column bottom */}
          <div className="bottom-grid">
            <ConditionBreakdown data={MOCK_BREAKDOWN} />
            <PeakHoursCard hours={PEAK_HOURS} />
          </div>

          {/* Session table full-width */}
          <SessionHistoryTable sessions={MOCK_SESSIONS} />
        </div>
      </main>

      {/* ── Bottom Nav (mobile) ── */}
      <nav className="bottom-nav">
        {navItems.map((item) => (
          <button
            key={item.id}
            className={`bottom-nav-item ${activeNav === item.id ? "bottom-nav-active" : ""}`}
            onClick={() => setActiveNav(item.id)}
          >
            <NavIcon id={item.id} />
            <span>{item.label}</span>
          </button>
        ))}
      </nav>

      <style jsx>{`
        * {
          box-sizing: border-box;
          margin: 0;
          padding: 0;
        }

        .root {
          display: flex;
          min-height: 100vh;
          background: #fafaf8;
          font-family: "DM Sans", "Helvetica Neue", Arial, sans-serif;
          color: #1a1a18;
        }

        /* ── Sidebar ── */
        .sidebar {
          width: 220px;
          flex-shrink: 0;
          background: #fff;
          border-right: 1.5px solid #e8e8e2;
          display: flex;
          flex-direction: column;
          padding: 24px 16px;
        }
        .sidebar-logo {
          padding: 4px 8px 28px;
        }
        .sidebar-nav {
          display: flex;
          flex-direction: column;
          gap: 2px;
          flex: 1;
        }

        .nav-item {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 12px;
          border-radius: 10px;
          border: none;
          background: transparent;
          cursor: pointer;
          font-size: 0.88rem;
          font-weight: 500;
          font-family: inherit;
          color: #888880;
          transition:
            background 0.15s,
            color 0.15s;
          text-align: left;
          width: 100%;
        }
        .nav-item:hover {
          background: #f5f5f0;
          color: #1a1a18;
        }
        .nav-active {
          background: #fef3c7 !important;
          color: #92400e !important;
        }

        .sidebar-footer {
          padding-top: 16px;
          border-top: 1.5px solid #e8e8e2;
        }
        .logout-btn {
          display: flex;
          align-items: center;
          gap: 10px;
          padding: 10px 12px;
          border-radius: 10px;
          border: none;
          background: transparent;
          cursor: pointer;
          font-size: 0.85rem;
          font-weight: 500;
          font-family: inherit;
          color: #888880;
          transition:
            background 0.15s,
            color 0.15s;
          width: 100%;
        }
        .logout-btn:hover {
          background: #fdf0ef;
          color: #c0392b;
        }

        /* ── Main ── */
        .main {
          flex: 1;
          min-width: 0;
          display: flex;
          flex-direction: column;
        }

        .topbar {
          display: flex;
          align-items: center;
          gap: 14px;
          padding: 20px 28px 0;
        }
        .topbar-logo-mobile {
          display: none;
        }
        .page-title {
          font-size: 1.4rem;
          font-weight: 700;
          letter-spacing: -0.025em;
          color: #1a1a18;
        }
        .page-subtitle {
          font-size: 0.82rem;
          color: #888880;
          margin-top: 2px;
        }

        /* Range toggle */
        .range-toggle {
          display: flex;
          background: #f0f0ea;
          border-radius: 10px;
          padding: 3px;
          gap: 2px;
          margin-left: auto;
        }
        .toggle-btn {
          padding: 6px 14px;
          border-radius: 8px;
          border: none;
          background: transparent;
          cursor: pointer;
          font-size: 0.8rem;
          font-weight: 500;
          font-family: inherit;
          color: #888880;
          transition:
            background 0.15s,
            color 0.15s;
        }
        .toggle-btn:hover {
          color: #1a1a18;
        }
        .toggle-active {
          background: #fff !important;
          color: #1a1a18 !important;
          box-shadow: 0 1px 3px #00000010;
        }

        /* ── Content ── */
        .content {
          padding: 20px 28px 40px;
          flex: 1;
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        /* Summary strip */
        .summary-strip {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: 12px;
        }
        .strip-card {
          background: #fff;
          border: 1.5px solid #e8e8e2;
          border-radius: 14px;
          padding: 16px 18px;
        }
        .strip-label {
          font-size: 0.72rem;
          font-weight: 500;
          color: #888880;
          text-transform: uppercase;
          letter-spacing: 0.06em;
          margin-bottom: 6px;
        }
        .strip-value {
          font-size: 1.4rem;
          font-weight: 700;
          letter-spacing: -0.03em;
          color: #1a1a18;
        }

        /* Bottom two-column */
        .bottom-grid {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
        }

        /* ── Bottom Nav ── */
        .bottom-nav {
          display: none;
        }

        /* ── Responsive ── */
        @media (max-width: 768px) {
          .sidebar {
            display: none;
          }
          .topbar-logo-mobile {
            display: flex;
          }
          .topbar {
            padding: 16px 18px 14px;
            border-bottom: 1.5px solid #e8e8e2;
            flex-wrap: wrap;
          }
          .range-toggle {
            margin-left: 0;
            margin-top: 8px;
            width: 100%;
            justify-content: stretch;
          }
          .toggle-btn {
            flex: 1;
          }
          .content {
            padding: 16px 18px 90px;
          }
          .summary-strip {
            grid-template-columns: 1fr 1fr;
          }
          .bottom-grid {
            grid-template-columns: 1fr;
          }

          .bottom-nav {
            display: flex;
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            background: #fff;
            border-top: 1.5px solid #e8e8e2;
            padding: 8px 0 12px;
            z-index: 100;
          }
          .bottom-nav-item {
            flex: 1;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 4px;
            border: none;
            background: transparent;
            cursor: pointer;
            font-size: 0.68rem;
            font-weight: 500;
            font-family: inherit;
            color: #888880;
            padding: 6px 0;
            transition: color 0.15s;
          }
          .bottom-nav-active {
            color: #b45309 !important;
          }
        }
      `}</style>
    </div>
  );
}
