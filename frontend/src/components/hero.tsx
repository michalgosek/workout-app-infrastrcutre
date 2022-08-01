const HERO_LOGO = 'https://cdn.vectorstock.com/i/1000x1000/52/08/fitness-club-designs-with-exercising-athletic-man-vector-25825208.webp';

const Hero = () => {
    return (
        <div className="text-center">
            <p className="lead">
                What hurts today makes you stronger tomorrow.
            </p>
            <p className="text-secondary">~ Jay Cutler 4x times Mr. Olympia </p>
            <img className="mb-2 app-logo" src={HERO_LOGO} alt="React logo" width="400" />
        </div>
    )
};

export default Hero;
