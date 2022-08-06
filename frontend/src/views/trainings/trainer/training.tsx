import { Container, Table } from 'react-bootstrap';
import { FC, PropsWithChildren, useEffect, useState } from 'react';
import { ParticipantReadModel, TrainerGroupReadModel } from 'services/models';

import DropButton from 'components/buttons/drop-button';
import TrainingForm from '../../../components/forms/form';
import { TrainingsService } from 'services/trainings-service';
import style from './training.module.scss';
import { useParams } from 'react-router-dom';

type ParticipantsTableProps = {
    trainingUUID: string;
    trainerUUID: string;
    participants: ParticipantReadModel[];
};

type deleteTrainingData = {
    trainerUUID?: string;
    particiapntUUID?: string;
    trainingUUID?: string;
};

const deleteParticipant = (data: deleteTrainingData) => {
    const params = [
        data.particiapntUUID ?? '',
        data.trainerUUID ?? '',
        data.trainingUUID ?? ''
    ];
    const isMissingValue = params.some(v => v === "");
    if (isMissingValue) {
        console.error('missing params value');
        return
    }
    TrainingsService.cancelParticipantWorkout(params[0], params[1], params[2])
        .then((res) => console.log(res))
        .catch(err => console.error(err));
};

const ParticipantsTable: FC<PropsWithChildren<ParticipantsTableProps>> = ({ participants, trainerUUID, trainingUUID }) => {
    return (
        <div className={style.rowtext}>
            <Table striped bordered hover>
                <thead className={style['t-headers']}>
                    <tr>
                        <th>#</th>
                        <th>Participant Name</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        participants.map((p, idx) => (
                            <tr key={p.uuid} className={style['t-rows']} >
                                <td>{idx + 1}</td>
                                <td>{p.name}</td>
                                <td>
                                    <DropButton onClickHandle={() => deleteParticipant({
                                        particiapntUUID: p.uuid,
                                        trainerUUID: trainerUUID,
                                        trainingUUID: trainingUUID,
                                    })} />
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};

const NoTrainerGroupParticipants: FC = () => {
    return <h1>No participants.</h1>
};

const NoTrainerGroup: FC = () => {
    return <h1>Trainer Group Not Exist.</h1>
};

const TrainerGroup: FC = () => {
    const { trainingUUID, trainerUUID } = useParams();
    const [training, setTraining] = useState<TrainerGroupReadModel>();

    useEffect(() => {
        if (!trainerUUID || !trainingUUID) {
            console.error('missing trainerUUID or trainingUUID in URL path');
            return;
        }
        TrainingsService.getTrainerGroup(trainerUUID, trainingUUID)
            .then(t => setTraining(t))
            .catch(err => console.error(err));
    }, [trainerUUID, trainingUUID]);

    if (!training || !trainerUUID || !trainingUUID) return <NoTrainerGroup />;

    return (
        <>
            <Container>
                <TrainingForm placeholders={{
                    name: training.name,
                    description: training.description,
                    appointment: training.date,
                }}
                    callbackPutAPI={TrainingsService.updateTrainingGroup} />
            </Container>

            <Container className='mt-3'>
                {!training.participants ? <NoTrainerGroupParticipants /> : <ParticipantsTable participants={training.participants} trainerUUID={trainerUUID} trainingUUID={trainingUUID} />}
            </Container>
        </>
    );
};


export default TrainerGroup;


