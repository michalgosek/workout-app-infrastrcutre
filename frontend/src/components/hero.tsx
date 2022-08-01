import logo from '../assets/hero.jpg';

const Hero = () => {
    return (
        <div className="text-center">
            <p className="lead">
                 "What hurts today makes you stronger tomorrow..." 🏋️‍
            </p>
            <img className="mb-2 app-logo" src={logo} alt="React logo" width="400" />
        </div>
    )
};

export default Hero;
