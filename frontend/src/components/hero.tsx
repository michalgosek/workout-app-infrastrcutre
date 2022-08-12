import { Container, Image } from 'react-bootstrap';

import logo from 'assets/hero.jpg';

const Hero = () => {
    return (
        <>
            <Container className='text-center'>
                <p className="lead">
                    "What hurts today makes you stronger tomorrow..." 🏋️‍
                </p>
                <Image className="mb-2 app-logo" src={logo} alt="React logo" width="400" />
            </Container>
        </>
    )
};

export default Hero;
