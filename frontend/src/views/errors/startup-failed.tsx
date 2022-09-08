import Alert from 'react-bootstrap/Alert';
import { Container } from 'react-bootstrap';

const StartupFailed = () => {
    return (
        <Container className='text-center'>
            <Alert variant="danger mt-3">
                <Alert.Heading>Failed to load .env app file</Alert.Heading>
                <p>
                    Please ensure that .env file has been created inside the repository
                    before starting the app.
                </p>
                <hr />
                <p className="mb-0">
                    Example structure of .env file:
                </p>
                <hr />
                <div>
                    <p>REACT_APP_AUTH0_DOMAIN=dev.auth0.com</p>
                    <p>REACT_APP_AUTH0_CLIENT_ID=A2ZimOBXMWfMlShbykLZeg1GI1GC2o0O</p>
                    <p>REACT_APP_AUTH0_AUDIENCE=https://service-name</p>
                    <p>REACT_APP_AUTH0_CALLBACK_URL=http://localhost:3000</p>
                </div>
            </Alert>
        </Container>
    );
};

export default StartupFailed;
