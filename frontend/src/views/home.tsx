import { Fragment } from 'react';
import Hero from '../components/hero';
import PageLayout from './page-layout';

const Home: React.FC = () => {
    return (
        <PageLayout>
            <Fragment>
                <Hero />
            </Fragment>
        </PageLayout>
    );
};

export default Home;
