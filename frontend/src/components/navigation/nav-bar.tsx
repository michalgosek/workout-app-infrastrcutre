import AuthNav from './auth-nav';
import MainNav from './main-nav';

const NavBar: React.FC = () => {
    return (
        <div className="nav-container mb-3">
            <nav className="navbar navbar-expand-md navbar-light bg-light">
                <div className="container">
                    <div className="navbar-brand" />
                    <MainNav />
                    <AuthNav />
                </div>
            </nav>
        </div>
    );
};

export default NavBar;