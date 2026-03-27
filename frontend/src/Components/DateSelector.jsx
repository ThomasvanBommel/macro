import DateString from '../Classes/DateString';

/**
 * Callback to set the selected date. Should accept a DateString object.
 * @callback setDate
 * @param { DateString } newDate - The new selected date
 * @returns { void }
 */

/**
 * A simple date selector component with "Prev" and "Next" buttons to increment/decrement the date,
 * and an input field to select a specific date.
 * @param { Object } props - The component props
 * @param { DateString } props.date - The currently selected date
 * @param { setDate } props.setDate - Function to update the selected date
 */
export default function DateSelector({ date, setDate }) {
    const incr = () => setDate(date.increment());
    const decr = () => setDate(date.decrement());

    return (
        <form onSubmit={ e => e.preventDefault() }>
            <fieldset role="group">
                <button onClick={ decr }>Prev</button>
                <input type="date" 
                       value={ date.value } 
                       onChange={ e => setDate(new DateString(e.target.value)) } />
                <button onClick={ incr }>Next</button>
            </fieldset>
        </form>
    );
}