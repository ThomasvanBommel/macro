export default class DateString {

    /**
     * Creates a new DateString instance.
     * @param { string } datestr - A date string in the format "YYYY-MM-DD".
     */
    constructor(datestr) {
        this.value = datestr;
    }

    /**
     * Converts the DateString to a Date object.
     * @return { Date } A Date object representing the date string.
     */
    toDate() {
        return new Date(this.value + " 00:00:00");
    }

    /**
     * Adds a specified number of days to the DateString.
     * @param { number } n - The number of days to add.
     * @return { DateString } A new DateString instance with the updated date.
     */
    addDays(n) {
        const date = this.toDate();
        date.setDate(date.getDate() + n);
        return new DateString(DateString.#dateToString(date));
    }

    /**
     * Increments the DateString by a specified number of days.
     * @param { number } [n=1] - The number of days to increment.
     * @return { DateString } A new DateString instance with the updated date.
     */
    increment(n=1) {
        return this.addDays(n);
    }

    /**
     * Decrements the DateString by a specified number of days.
     * @param { number } [n=1] - The number of days to decrement.
     * @return { DateString } A new DateString instance with the updated date.
     */
    decrement(n=1) {
        return this.addDays(-n);
    }

    /**
     * Returns a DateString instance representing today's date.
     * @return { DateString } A DateString instance for today's date.
     */
    static today() {
        return new DateString(DateString.#dateToString(new Date()));
    }

    /**
     * Converts a Date object to a string in the format "YYYY-MM-DD".
     * @param { Date } date - A Date object to be converted.
     * @return { string } A string representing the date in the format "YYYY-MM-DD".
     */
    static #dateToString(date) {
        return date.toLocaleString().split(',')[0];
    }
}