"use client";

import { useState } from "react";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    // TODO: Replace with actual API call
    try {
      await new Promise((resolve) => setTimeout(resolve, 1200));
      // await fetch('/api/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) })
      console.log("Login submitted", { email, password });
    } catch {
      setError("Invalid credentials. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-root">
      <div className="login-left">
        <div className="brand">
          <BulbEyeLogo />
          <span className="brand-name">StareDesk</span>
          <span className="brand-tagline">Sit, Stay &amp; Focus</span>
        </div>
      </div>

      <div className="login-right">
        <form className="login-form" onSubmit={handleSubmit} noValidate>
          <div className="form-header">
            <h1>Welcome back</h1>
            <p>Sign in to your workspace</p>
          </div>

          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={loading}
            />
          </div>

          <div className="field">
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={loading}
            />
          </div>

          {error && <p className="error-msg">{error}</p>}

          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? <span className="spinner" /> : "Sign in"}
          </button>
        </form>
      </div>

      <style jsx>{`
        /* ── Reset & Root ── */
        .login-root {
          min-height: 100vh;
          display: flex;
          font-family: "DM Sans", "Helvetica Neue", Arial, sans-serif;
          background: #fafaf8;
          color: #1a1a1a;
        }

        /* ── Left Panel ── */
        .login-left {
          flex: 1;
          display: flex;
          align-items: center;
          justify-content: center;
          background: #1a1a18;
          position: relative;
          overflow: hidden;
        }

        .login-left::before {
          content: "";
          position: absolute;
          width: 520px;
          height: 520px;
          border-radius: 50%;
          background: radial-gradient(circle, #f5a82320 0%, transparent 70%);
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          pointer-events: none;
        }

        .brand {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 14px;
          z-index: 1;
        }

        .brand-name {
          font-size: 2.5rem;
          font-weight: 700;
          letter-spacing: -0.03em;
          color: #f5f5f0;
        }

        .brand-tagline {
          font-size: 0.95rem;
          font-weight: 400;
          color: #888880;
          letter-spacing: 0.01em;
        }

        /* ── Right Panel ── */
        .login-right {
          width: 480px;
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 3rem;
          background: #fafaf8;
        }

        .login-form {
          width: 100%;
          max-width: 360px;
          display: flex;
          flex-direction: column;
          gap: 1.5rem;
        }

        .form-header h1 {
          font-size: 1.75rem;
          font-weight: 700;
          letter-spacing: -0.025em;
          color: #1a1a18;
          margin: 0 0 6px;
        }

        .form-header p {
          font-size: 0.9rem;
          color: #888880;
          margin: 0;
        }

        /* ── Fields ── */
        .field {
          display: flex;
          flex-direction: column;
          gap: 6px;
        }

        label {
          font-size: 0.82rem;
          font-weight: 500;
          color: #555550;
          letter-spacing: 0.02em;
          text-transform: uppercase;
        }

        input {
          width: 100%;
          box-sizing: border-box;
          padding: 12px 14px;
          font-size: 0.95rem;
          font-family: inherit;
          color: #1a1a18;
          background: #fff;
          border: 1.5px solid #e4e4de;
          border-radius: 10px;
          outline: none;
          transition:
            border-color 0.18s,
            box-shadow 0.18s;
          -webkit-appearance: none;
        }

        input::placeholder {
          color: #bbbbb5;
        }

        input:hover:not(:disabled) {
          border-color: #c8c8c0;
        }

        input:focus {
          border-color: #f5a823;
          box-shadow: 0 0 0 3px #f5a82322;
        }

        input:disabled {
          opacity: 0.55;
          cursor: not-allowed;
        }

        /* ── Error ── */
        .error-msg {
          font-size: 0.85rem;
          color: #c0392b;
          margin: -8px 0 0;
          padding: 10px 12px;
          background: #fdf0ef;
          border-radius: 8px;
          border: 1px solid #f5c6c4;
        }

        /* ── Submit Button ── */
        .submit-btn {
          width: 100%;
          padding: 13px;
          font-size: 0.95rem;
          font-weight: 600;
          font-family: inherit;
          color: #1a1a18;
          background: #f5a823;
          border: none;
          border-radius: 10px;
          cursor: pointer;
          display: flex;
          align-items: center;
          justify-content: center;
          transition:
            background 0.18s,
            transform 0.1s;
          letter-spacing: -0.01em;
        }

        .submit-btn:hover:not(:disabled) {
          background: #e09618;
        }

        .submit-btn:active:not(:disabled) {
          transform: scale(0.985);
        }

        .submit-btn:disabled {
          cursor: not-allowed;
          opacity: 0.7;
        }

        /* ── Spinner ── */
        .spinner {
          display: inline-block;
          width: 18px;
          height: 18px;
          border: 2.5px solid #1a1a1840;
          border-top-color: #1a1a18;
          border-radius: 50%;
          animation: spin 0.7s linear infinite;
        }

        @keyframes spin {
          to {
            transform: rotate(360deg);
          }
        }

        /* ── Responsive: Mobile ── */
        @media (max-width: 768px) {
          .login-root {
            flex-direction: column;
          }

          .login-left {
            flex: 0;
            padding: 3rem 2rem 2.5rem;
          }

          .login-left::before {
            display: none;
          }

          .brand {
            flex-direction: row;
            gap: 14px;
            align-items: center;
          }

          .brand-name {
            font-size: 1.8rem;
          }

          .brand-tagline {
            display: none;
          }

          .login-right {
            width: 100%;
            padding: 2.5rem 1.5rem 3rem;
            flex: 1;
            align-items: flex-start;
          }

          .login-form {
            max-width: 100%;
          }
        }
      `}</style>
    </div>
  );
}

function BulbEyeLogo() {
  return (
    <svg
      width="56"
      height="68"
      viewBox="0 0 56 68"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-label="StareDesk logo"
    >
      {/* Bulb body */}
      <path
        d="M28 4C15.85 4 6 13.85 6 26c0 8.4 4.6 15.7 11.4 19.6V52h21.2v-6.4C45.4 41.7 50 34.4 50 26 50 13.85 40.15 4 28 4z"
        fill="#F5A823"
      />
      {/* Base / cap */}
      <rect x="17.4" y="54" width="21.2" height="8" rx="4" fill="#F5A823" />

      {/* Eye white */}
      <ellipse cx="28" cy="25" rx="10" ry="6.5" fill="white" />
      {/* Pupil */}
      <circle cx="28" cy="25" r="4" fill="#1a1a18" />
      {/* Pupil highlight */}
      <circle cx="29.6" cy="23.4" r="1.2" fill="white" />

      {/* Lash lines — top */}
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
      {/* Lash lines — bottom */}
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
  );
}
