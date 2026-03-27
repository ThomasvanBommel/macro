import { useState, useEffect } from 'react';

/**
 * Called when the date changes, either through the input or the buttons. 
 * @callback onDateChange
 * @param { string } datestr - the new date
 */

/**
 * Component to increment, decrement, or set a local date string.
 * @param { string } [datestr] - the date to initially show, in YYYY-MM-DD format. If not provided, 
 *                               defaults to the current date.
 * @param { onDateChange } [onDateChange] - callback for when the date changes 
 * @returns DateSelector jsx component
 */
export default function DateSelector({ datestr, onDateChange }) {
    const [date, setDate] = useState(datestr ?? today());

    const increment = () => setDate(addDaysToDateString(date, 1));
    const decrement = () => setDate(addDaysToDateString(date, -1));

    useEffect(() => {
        onDateChange?.(date);
    }, [date]);

    return (
        <form onSubmit={ e => e.preventDefault() }>
            <fieldset role="group">
                <button onClick={ decrement }>Prev</button>
                <input type="date" value={ date } onChange={ e => setDate(e.target.value) } />
                <button onClick={ increment }>Next</button>
            </fieldset>
        </form>
    );
}

/**
 * Gets the current date in YYYY-MM-DD format.
 * @returns { string } - the current date in YYYY-MM-DD format
 */
function today() {
    return dateToString(new Date());
}

/**
 * Converts a Date object to a string in YYYY-MM-DD format.
 * @param { Date } date - the date to convert
 * @returns { string } - the date in YYYY-MM-DD format
 */
function dateToString(date) {
    return date.toISOString().split('T')[0];
}

/**
 * Adds a number of days to a date string.
 * @param { string } datestr - the date string in YYYY-MM-DD format
 * @param { integer } n - the number of days to add (can be negative)
 * @returns { string } - the new date string in YYYY-MM-DD format
 */
function addDaysToDateString(datestr, n) {
    const date = new Date(datestr);
    date.setDate(date.getDate() + n);
    return dateToString(date);
}