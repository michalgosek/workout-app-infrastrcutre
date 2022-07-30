import PageLayout from './page-layout';
import { useAuth0 } from '@auth0/auth0-react';

const Profile: React.FC = () => {
    const { user } = useAuth0();
    if (!user) {
        return null
    }
    return (
        <PageLayout>
            <div className="row algin-items-center profile-header">
                <div className="col-md-3 mb-3">
                    <img
                        src={user.picture}
                        alt="Profile"
                        className="rounded-circle img-fluid profile-picture mb-3 mb-md-0" />
                </div>
                <div className="col-md text-center text-md-left">
                    <h2>{user.name}</h2>
                    <p className="lead text-muted">{user.email}</p>
                </div>
            </div>
        </PageLayout>
    );
};

export default Profile;
