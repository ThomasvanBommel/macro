
// Home.jsx
export default function Home() {
    return (
        <div>
            <header>
                <hgroup>
                    <p>Portfolio build • Full-stack nutrition tracker</p>
                </hgroup>
                <p>
                    Macro demonstrates practical end-to-end delivery: schema and API design,
                    session-backed auth flows, and a working React UI for daily macro tracking.
                </p>
                <p>
                    Source code: <a href="https://github.com/ThomasvanBommel/macro" target="_blank" rel="noopener noreferrer">GitHub</a>
                </p>
            </header>

            <section>
                <article>
                    <h3>🧱 Stack Snapshot</h3>
                    <ul>
                        <li>React 19 + Vite 8</li>
                        <li>Go 1.26 + Gin</li>
                        <li>SQLite + Goose migrations</li>
                        <li>Docker + Compose + Make</li>
                        <li>Cookie-backed sessions</li>
                    </ul>
                    <p>
                        Architecture: <code>React (Vite) -&gt; /api -&gt; Go (Gin) -&gt; SQLite</code>
                    </p>
                </article>
            </section>

            <section>
                <article>
                    <h3>Why this exists</h3>
                    <p>
                        This project is built as a realistic vertical slice: shipping useful
                        functionality first, then iterating on quality, UX, and operational
                        hardening.
                    </p>
                </article>
            </section>

            <section className="grid">
                <article>
                    <h3>✅ Done</h3>
                    <ul>
                        <li>Auth flow and cookie-backed sessions</li>
                        <li>Food search and custom food creation</li>
                        <li>Date-based entries with daily macro totals</li>
                    </ul>
                </article>

                <article>
                    <h3>🚧 In Progress</h3>
                    <ul>
                        <li>Frontend tests for API helpers and key flows</li>
                        <li>Centralized user-friendly API error mapping</li>
                        <li>Profile UX and data presentation improvements</li>
                    </ul>
                </article>

                <article>
                    <h3>🧭 Next Up</h3>
                    <ul>
                        <li>Edit and remove diary entries</li>
                        <li>Edit food records and show macro percentages</li>
                        <li>Deployment pipeline automation</li>
                    </ul>
                </article>
            </section>

            <details>
                <summary>Reviewer context</summary>
                This is intentionally a working, evolving product rather than a static demo.
                The goal is to show ownership across product, backend, frontend, and delivery.
            </details>
        </div>
    );
}