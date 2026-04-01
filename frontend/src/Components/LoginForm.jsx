import { handleLoginUser } from '../api'

/**
 * Login component that renders a login form and handles user authentication. It takes an onSuccess
 * callback prop that is called with the session data upon successful login. The form collects the
 * username and password from the user and submits it to the server using the handleLoginUser 
 * function. If the login is successful, the onSuccess callback is called with the session data. If 
 * the login fails, an error message is displayed to the user.
 * @param {Object} props - The component props.
 * @param {Function} props.onSuccess - The callback function to be called upon successful login.
 * @returns {JSX.Element} The rendered Login component.
 */
export default function LoginForm({ onSuccess }) {
    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);
        
        handleLoginUser(form.get('name'), form.get('password'))
            .then(onSuccess)
            .catch(err => alert(`Login failed: ${err.message}`));
    };

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Username:
                <input type="text" id="name" name="name" required />
            </label>
            <label>
                Password:
                <input type="password" id="password" name="password" required />
            </label>
            <button type="submit">Login</button>
        </form>
    );
}