"use client";

import { useState, useEffect } from "react";

// ── Types ──────────────────────────────────────────────────────────────────
type Condition =
  | "OPTIMAL"
  | "EYE_STRAIN_RISK"
  | "POSTURE_RISK"
  | "DISTRACTED"
  | "AWAY";
type NavItem = "dashboard" | "analytics" | "settings";

interface SensorData {
  condition: Condition;
  distance_cm: number;
  ldr_value: number;
  pir_detected: boolean;
}

interface TodaySummary {
  focus_minutes: number;
  peak_hour: string;
  sessions_count: number;
}

// ── Constants ──────────────────────────────────────────────────────────────
const CONDITION_META: Record<
  Condition,
  { label: string; color: string; bg: string; dot: string }
> = {
  OPTIMAL: {
    label: "Optimal",
    color: "#1D6A3A",
    bg: "#E6F4EC",
    dot: "#22C55E",
  },
  EYE_STRAIN_RISK: {
    label: "Eye Strain Risk",
    color: "#92400E",
    bg: "#FEF3C7",
    dot: "#F59E0B",
  },
  POSTURE_RISK: {
    label: "Posture Risk",
    color: "#92400E",
    bg: "#FEF3C7",
    dot: "#F59E0B",
  },
  DISTRACTED: {
    label: "Distracted",
    color: "#9A3412",
    bg: "#FFF0E6",
    dot: "#F97316",
  },
  AWAY: { label: "Away", color: "#374151", bg: "#F3F4F6", dot: "#9CA3AF" },
};

// ── Mock data (replace with WebSocket / API) ───────────────────────────────
const MOCK_SENSOR: SensorData = {
  condition: "OPTIMAL",
  distance_cm: 62,
  ldr_value: 718,
  pir_detected: true,
};

const MOCK_SUMMARY: TodaySummary = {
  focus_minutes: 142,
  peak_hour: "09:00 – 10:00",
  sessions_count: 4,
};

// ── Sub-components ─────────────────────────────────────────────────────────

function Logo({ small = false }: { small?: boolean }) {
  return (
    <div style={{ display: "flex", alignItems: "center", gap: small ? 8 : 10 }}>
      <svg
        width={small ? 28 : 34}
        height={small ? 34 : 42}
        viewBox="0 0 56 68"
        fill="none"
      >
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
      {!small && (
        <span
          style={{
            fontSize: "1.1rem",
            fontWeight: 700,
            letterSpacing: "-0.03em",
            color: "#1a1a18",
          }}
        >
          StareDesk
        </span>
      )}
    </div>
  );
}

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

function DeviceStatusBar({
  isOnline,
  configPending,
}: {
  isOnline: boolean;
  configPending: boolean;
}) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: 12,
        padding: "10px 16px",
        background: "#fff",
        border: "1.5px solid #e8e8e2",
        borderRadius: 12,
        marginBottom: 20,
      }}
    >
      <span
        style={{
          width: 8,
          height: 8,
          borderRadius: "50%",
          background: isOnline ? "#22C55E" : "#9CA3AF",
          flexShrink: 0,
        }}
      />
      <span
        style={{
          fontSize: "0.85rem",
          fontWeight: 500,
          color: isOnline ? "#1D6A3A" : "#6B7280",
        }}
      >
        Device {isOnline ? "online" : "offline"}
      </span>
      {configPending && (
        <span
          style={{
            marginLeft: "auto",
            fontSize: "0.75rem",
            fontWeight: 500,
            padding: "3px 10px",
            borderRadius: 20,
            background: "#FEF3C7",
            color: "#92400E",
          }}
        >
          Config pending
        </span>
      )}
      {!configPending && isOnline && (
        <span
          style={{
            marginLeft: "auto",
            fontSize: "0.75rem",
            fontWeight: 500,
            padding: "3px 10px",
            borderRadius: 20,
            background: "#E6F4EC",
            color: "#1D6A3A",
          }}
        >
          Config synced
        </span>
      )}
    </div>
  );
}

function LiveConditionCard({ sensor }: { sensor: SensorData }) {
  const meta = CONDITION_META[sensor.condition];
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
          margin: "0 0 14px",
        }}
      >
        Live condition
      </p>
      <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
        <span
          style={{
            width: 14,
            height: 14,
            borderRadius: "50%",
            background: meta.dot,
            flexShrink: 0,
          }}
        />
        <span
          style={{
            fontSize: "1.6rem",
            fontWeight: 700,
            letterSpacing: "-0.025em",
            color: "#1a1a18",
          }}
        >
          {meta.label}
        </span>
        <span
          style={{
            marginLeft: "auto",
            fontSize: "0.78rem",
            fontWeight: 500,
            padding: "4px 12px",
            borderRadius: 20,
            background: meta.bg,
            color: meta.color,
          }}
        >
          {sensor.pir_detected ? "Present" : "No presence"}
        </span>
      </div>
    </div>
  );
}

function SessionTimer() {
  const [seconds, setSeconds] = useState(2740); // mock: 45m 40s

  useEffect(() => {
    const id = setInterval(() => setSeconds((s) => s + 1), 1000);
    return () => clearInterval(id);
  }, []);

  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  const fmt = (n: number) => String(n).padStart(2, "0");

  return (
    <div
      style={{
        background: "#1a1a18",
        border: "1.5px solid #1a1a18",
        borderRadius: 16,
        padding: "20px 22px",
        display: "flex",
        flexDirection: "column",
        gap: 10,
      }}
    >
      <p
        style={{
          fontSize: "0.75rem",
          fontWeight: 500,
          color: "#888880",
          textTransform: "uppercase",
          letterSpacing: "0.06em",
          margin: 0,
        }}
      >
        Session timer
      </p>
      <p
        style={{
          fontSize: "2.4rem",
          fontWeight: 700,
          letterSpacing: "-0.04em",
          color: "#F5A823",
          margin: 0,
          fontVariantNumeric: "tabular-nums",
        }}
      >
        {h > 0 ? `${fmt(h)}:` : ""}
        {fmt(m)}:{fmt(s)}
      </p>
      <p style={{ fontSize: "0.8rem", color: "#555550", margin: 0 }}>
        Started at 09:14
      </p>
    </div>
  );
}

function SensorReadings({ sensor }: { sensor: SensorData }) {
  const items = [
    {
      label: "Distance",
      value: `${sensor.distance_cm} cm`,
      sub: "from monitor",
    },
    { label: "Light level", value: sensor.ldr_value, sub: "LDR value" },
    {
      label: "Presence",
      value: sensor.pir_detected ? "Detected" : "None",
      sub: "PIR sensor",
    },
  ];

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
          margin: "0 0 16px",
        }}
      >
        Sensor readings
      </p>
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(3, 1fr)",
          gap: 12,
        }}
      >
        {items.map((item) => (
          <div
            key={item.label}
            style={{
              background: "#fafaf8",
              border: "1px solid #e8e8e2",
              borderRadius: 10,
              padding: "12px 14px",
            }}
          >
            <p
              style={{
                fontSize: "0.72rem",
                color: "#888880",
                margin: "0 0 4px",
                fontWeight: 500,
              }}
            >
              {item.label}
            </p>
            <p
              style={{
                fontSize: "1.1rem",
                fontWeight: 700,
                color: "#1a1a18",
                margin: "0 0 2px",
                letterSpacing: "-0.02em",
              }}
            >
              {item.value}
            </p>
            <p style={{ fontSize: "0.7rem", color: "#aaaaaa", margin: 0 }}>
              {item.sub}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}

function TodaySummaryCard({ summary }: { summary: TodaySummary }) {
  const hours = Math.floor(summary.focus_minutes / 60);
  const mins = summary.focus_minutes % 60;
  const label = hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;

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
          margin: "0 0 16px",
        }}
      >
        Today&apos;s summary
      </p>
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
        <div
          style={{
            background: "#fafaf8",
            border: "1px solid #e8e8e2",
            borderRadius: 10,
            padding: "14px 16px",
          }}
        >
          <p
            style={{
              fontSize: "0.72rem",
              color: "#888880",
              margin: "0 0 4px",
              fontWeight: 500,
            }}
          >
            Focus time
          </p>
          <p
            style={{
              fontSize: "1.5rem",
              fontWeight: 700,
              color: "#1a1a18",
              margin: 0,
              letterSpacing: "-0.03em",
            }}
          >
            {label}
          </p>
        </div>
        <div
          style={{
            background: "#fafaf8",
            border: "1px solid #e8e8e2",
            borderRadius: 10,
            padding: "14px 16px",
          }}
        >
          <p
            style={{
              fontSize: "0.72rem",
              color: "#888880",
              margin: "0 0 4px",
              fontWeight: 500,
            }}
          >
            Sessions
          </p>
          <p
            style={{
              fontSize: "1.5rem",
              fontWeight: 700,
              color: "#1a1a18",
              margin: 0,
              letterSpacing: "-0.03em",
            }}
          >
            {summary.sessions_count}
          </p>
        </div>
        <div
          style={{
            background: "#fafaf8",
            border: "1px solid #e8e8e2",
            borderRadius: 10,
            padding: "14px 16px",
            gridColumn: "span 2",
          }}
        >
          <p
            style={{
              fontSize: "0.72rem",
              color: "#888880",
              margin: "0 0 4px",
              fontWeight: 500,
            }}
          >
            Peak hour
          </p>
          <p
            style={{
              fontSize: "1rem",
              fontWeight: 700,
              color: "#1a1a18",
              margin: 0,
              letterSpacing: "-0.02em",
            }}
          >
            {summary.peak_hour}
          </p>
        </div>
      </div>
    </div>
  );
}

// ── Main Dashboard Page ────────────────────────────────────────────────────
export default function DashboardPage() {
  const [activeNav, setActiveNav] = useState<NavItem>("dashboard");
  const sensor = MOCK_SENSOR;
  const summary = MOCK_SUMMARY;

  const navItems: { id: NavItem; label: string }[] = [
    { id: "dashboard", label: "Dashboard" },
    { id: "analytics", label: "Analytics" },
    { id: "settings", label: "Settings" },
  ];

  return (
    <div className="root">
      {/* ── Sidebar (desktop) ── */}
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

      {/* ── Main Content ── */}
      <main className="main">
        <header className="topbar">
          <div className="topbar-logo-mobile">
            <Logo small />
          </div>
          <div>
            <h1 className="page-title">Dashboard</h1>
            <p className="page-subtitle">Monday, 28 April 2026</p>
          </div>
        </header>

        <div className="content">
          <DeviceStatusBar isOnline={true} configPending={false} />

          <div className="grid-main">
            <div className="col-left">
              <LiveConditionCard sensor={sensor} />
              <SensorReadings sensor={sensor} />
            </div>
            <div className="col-right">
              <SessionTimer />
              <TodaySummaryCard summary={summary} />
            </div>
          </div>
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

        .content {
          padding: 20px 28px 28px;
          flex: 1;
        }

        .grid-main {
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 16px;
        }

        .col-left,
        .col-right {
          display: flex;
          flex-direction: column;
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
            padding: 16px 18px 0;
            border-bottom: 1.5px solid #e8e8e2;
            padding-bottom: 14px;
          }

          .content {
            padding: 16px 18px 90px;
          }

          .grid-main {
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
