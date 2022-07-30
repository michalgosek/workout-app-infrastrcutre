const LOADING_IMAGE = 'https://cdn.auth0.com/blog/auth0-react-sample/assets/loading.svg';

const Loading = () => {
    return (
        <div className="spinner">
            <img src={LOADING_IMAGE} alt="Loading..." />
        </div>
    );
};

export default Loading;