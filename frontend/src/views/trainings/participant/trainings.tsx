import { Container } from 'react-bootstrap';
import {FC} from 'react';
import style from './trainings.module.scss';

const ParticipantTrainings: FC = () => {
    return (
        <Container>
        <h1 className={style.title}>There will be participant trainings</h1>
        </Container>
    );
};  

export default ParticipantTrainings; 