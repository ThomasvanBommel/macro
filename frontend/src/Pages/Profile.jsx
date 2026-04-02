import { useState, useContext } from 'react';
import { DateSelector, EntryList } from '../Components';
import DateString from '../Classes/DateString';
import { SessionContext } from '../Context';


export default function Profile() {
    const [date, setDate] = useState(DateString.today());
    const session = useContext(SessionContext);

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