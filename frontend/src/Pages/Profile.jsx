import { Diary } from '../Components';
import { useSession } from '../Context';

export default function Profile() {
    const { username } = useSession();

    return <Diary username={ username } />;
}