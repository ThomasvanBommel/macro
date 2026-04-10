import { useState, useContext } from 'react';
import { DateSelector, EntryList, Diary } from '../Components';
import DateString from '../Classes/DateString';
import { SessionContext } from '../Context';

export default function Profile() {
    const [date, setDate] = useState(DateString.today());
    const session = useContext(SessionContext);

    if (session.changingState)
        return <article aria-busy="true"></article>;

    return (
        <>  
            <Diary username={ session?.username } />

            {/* <div style={{ display: "grid", placeItems: "center" }}>
                <DateSelector date={ date } setDate={ setDate } />
            </div> */}

            {/* <EntryList username={ session?.username } date={ date } editable={ true } /> */}
        </>
    )
}