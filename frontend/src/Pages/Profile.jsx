import { useState } from 'react';
import { DateSelector, EntryList } from '../Components';
import DateString from '../Classes/DateString';


export default function Profile({ session }) {
    const [date, setDate] = useState(DateString.today());

    console.log(session);

    return (
        <>  
            <div>
                <h2>Profile Page</h2>
                <p>Welcome to your profile, {session?.username}!</p>
            </div>

            <div style={{ display: "grid", placeItems: "center" }}>
                <DateSelector date={ date } setDate={ setDate } />
            </div>

            <EntryList username={ session?.username } date={ date } editable={ true } />
        </>
    )
}