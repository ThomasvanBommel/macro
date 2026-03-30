import { useState, useEffect } from 'react';
import axios from 'axios';
import { fetchFoods } from '../api';
import { DateSelector, EntryList } from '../Components';
import DateString from '../Classes/DateString';

export default function Profile({ session }) {
    const [date, setDate] = useState(DateString.today());

    return (
        <>  
            <div>
                <h2>Profile Page</h2>
                <p>Welcome to your profile, {session?.user_name}!</p>
            </div>

            <div style={{ display: "grid", placeItems: "center" }}>
                <DateSelector date={ date } setDate={ setDate } />
            </div>

            <EntryList username={ session?.user_name } date={ date } editable={ true } />
        </>
    )
}