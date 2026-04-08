import DateString from '../Classes/DateString';

// Component to select a date, with buttons to increment and decrement the date, and an input to 
// select a specific date
export default function DateSelector({ date, setDate }) {
    const incr = () => setDate(date.increment());
    const decr = () => setDate(date.decrement());

    return (
        <form onSubmit={ e => e.preventDefault() }>
            <fieldset role="group">
                <button onClick={ decr }>&lt;</button>
                <input type="date" 
                       value={ date.value } 
                       onChange={ e => setDate(new DateString(e.target.value)) } />
                <button onClick={ incr }>&gt;</button>
            </fieldset>
        </form>
    );
}