
export default function Profile({ session }) {
    return (
        <div>
            <h2>Profile Page</h2>
            <p>Welcome to your profile, {session?.username}!</p>
            <p>Your user ID is: {session?.user_id}</p>
            <p>Session expires at: {session?.expires_at}</p>
        </div>
    )
}