
/**
 * This component serves as the homepage of the application, providing a welcome message and an 
 * overview of the app's features. It also includes a link to the GitHub repository where the source
 * code can be found. The component is designed to be simple and informative, giving users a clear 
 * understanding of what the app is about and how to get started.
 * @returns { JSX.Element } - A React component that renders the homepage content.
 */
export default function Home() {
    return (
        <div>
            <h2>Welcome!</h2>
            <p>
                This is a simple web app built with React and Go, designed to demonstrate my 
                fullstack potential as a developer.
            </p>
            <p>
                The app allows users to track their food intake and monitor their macronutrients. 
                You can register for an account, log in, and start adding foods and entries to see 
                your nutritional information.
            </p>
            <p>
                Check out the source code on <a href="https://github.com/ThomasvanBommel/macro" 
                                                target="_blank" 
                                                rel="noopener noreferrer">GitHub</a>.
            </p>
        </div>
    );
}