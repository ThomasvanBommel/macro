import { Link } from 'react-router';

export default function Home() {
    return (
        <div>
            <h2>Welcome!</h2>
            <p>
                This is a simple web app built with React and Go, designed to demonstrate my fullstack potential as a developer.
            </p>
            <p>
                Use the navigation links above to explore the app. You can <Link to="/register">register</Link> a new account, <Link to="/login">log in</Link>, and view your <Link to="/profile">profile</Link>.
            </p>
            <p>
                Feel free to check out the source code on <a href="https://github.com/your-repo/macro" target="_blank" rel="noopener noreferrer">GitHub</a>.
            </p>
        </div>
    );
}