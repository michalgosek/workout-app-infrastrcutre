const UNIVERSITY_LOGO = 'https://forumakademickie.pl/wp-content/uploads/2020/11/b0cea52e11208b91b8e595ea98cb4065.jpg';

const Footer = () => {
    return (<footer className="bg-light p-3 text-center">
        <p className="lead">
            Thesis: Design and implementation of web applications based on the microservices architecture developed on the example system for gym and fitness club management.
        </p>
        <p className="lead">Computer Science, Michał Gosek (S069973).</p>
        <img className="mb-2" src={UNIVERSITY_LOGO} alt="Kielce University of Technology" width="40" />
    </footer>
    );
};

export default Footer;
