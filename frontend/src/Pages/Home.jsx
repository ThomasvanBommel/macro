
// Home.jsx
export default function Home() {
    return (
        <div className="home-page">
            <article className="hero">
                <p className="eyebrow">Nutrition tracker</p>
                <h2>Track meals and macros with minimal friction.</h2>
                <p>
                    Macro helps you log what you eat, review your day by meal,
                    and keep an eye on calories, protein, carbs, and fat.
                </p>
            </article>

            <section className="home-grid">
                <article>
                    <h3>What you can do</h3>
                    <ul>
                        <li>Create foods and reuse them in entries</li>
                        <li>Log entries by date and meal</li>
                        <li>See daily totals and macro split</li>
                    </ul>
                </article>

                <article>
                    <h3>Quick start</h3>
                    <ol>
                        <li>Register or log in</li>
                        <li>Open your profile</li>
                        <li>Add your first entry</li>
                    </ol>
                </article>
            </section>

            <article>
                <h3>Project</h3>
                <p>
                    This app is built as a practical full-stack project with React, Go, and SQLite.
                </p>
                <p>
                    Source code: <a href="https://github.com/ThomasvanBommel/macro" target="_blank" rel="noopener noreferrer">GitHub</a>
                </p>
            </article>

            <style>
                {`
                    .home-page {
                        display: flex;
                        flex-direction: column;
                        gap: 1rem;
                    }

                    .home-page .hero {
                        border-left: 4px solid var(--pico-primary);
                        padding-left: 1rem;
                    }

                    .home-page .hero h2 {
                        margin: 0.25rem 0 0.75rem;
                    }

                    .home-page .eyebrow {
                        margin: 0;
                        text-transform: uppercase;
                        letter-spacing: 0.06em;
                        font-size: 0.8rem;
                        opacity: 0.7;
                    }

                    .home-grid {
                        display: grid;
                        grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
                        gap: 1rem;
                    }

                    .home-page article :is(ul, ol) {
                        margin-bottom: 0;
                    }
                `}
            </style>
        </div>
    );
}