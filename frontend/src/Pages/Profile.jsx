import { useState, useContext } from 'react';
import { Diary } from '../Components';
import DateString from '../Classes/DateString';
import { SessionContext } from '../Context';

export default function Profile() {
    const [date, setDate] = useState(DateString.today());
    const session = useContext(SessionContext);

    if (session.changingState)
        return <article aria-busy="true"></article>;

    return (
        <Diary username={ session?.username } />
    )
}