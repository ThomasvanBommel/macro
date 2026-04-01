import { handleRegisterUser } from '../api'

/**
 * Register component that renders a registration form and handles user registration. It takes an 
 * onSuccess callback that is called when the registration is successful. The form collects the 
 * username, password, and password confirmation from the user. It validates that the password and
 * confirmation match before submitting.
 * @param {Object} props - The component props.
 * @param {Function} props.onSuccess - Callback function to be called when registration is 
 *                                     successful.
 */
export default function RegisterForm({ onSuccess }) {
    const handleSubmit = e => {
        e.preventDefault();

        const form = new FormData(e.currentTarget);

        handleRegisterUser(form.get('name'), form.get('password'))
            .then(onSuccess)
            .catch(err => alert(`Registration failed: ${err.message}`));
    }

    const handlePasswordChange = e => {
        const target = e.currentTarget;
        const p1 = target.form.password;
        const p2 = target.form.confirm;

        target.setCustomValidity(p1.value !== p2.value ? "Passwords do not match" : "");
    }

    return (
        <form onSubmit={handleSubmit}>
            <label>
                Username:
                <input type="text" name="name" required />
            </label>
            <label>
                Password:
                <input type="password" name="password" onChange={handlePasswordChange} required />
            </label>
            <label>
                Confirm :
                <input type="password" name="confirm" onChange={handlePasswordChange} required />
            </label>
            <input type="submit" value="Register" />
        </form>
    );
}