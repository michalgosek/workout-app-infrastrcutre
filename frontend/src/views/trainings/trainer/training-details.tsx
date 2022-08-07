import { Container, Table } from 'react-bootstrap';
import { FC, PropsWithChildren, useEffect, useState } from 'react';
import { ParticipantReadModel, TrainerGroupReadModel } from 'services/models';

import DropButton from 'components/buttons/drop-button';
import TrainingForm from '../../../components/forms/form';
import { TrainingsService } from 'services/trainings-service';
import style from './training-details.module.scss';
import { useGetTrainingsServiceAccessToken } from 'services/hooks';
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

const deleteParticipant = (data: deleteTrainingData, token: string) => {
    const params = [
        data.particiapntUUID ?? '',
        data.trainerUUID ?? '',
        data.trainingUUID ?? '',
        token,
    ];
    const isMissingValue = params.some(v => v === "");
    if (isMissingValue) {
        console.error('missing params value');
        return
    }

    const cancelWorkout = async () => {
        try {
            const res = await TrainingsService.cancelParticipantWorkout(params[0], params[1], params[2], token);
            console.log(res);
            window.location.reload();
        } catch (error) {
            console.error(error)
        }
    }
    cancelWorkout();
};

const ParticipantsTable: FC<PropsWithChildren<ParticipantsTableProps>> = ({ participants, trainerUUID, trainingUUID }) => {
    const token = useGetTrainingsServiceAccessToken();
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
                                    }, token)} />
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
    const token = useGetTrainingsServiceAccessToken();

    useEffect(() => {
        if (!trainerUUID || !trainingUUID) {
            console.error('missing trainerUUID or trainingUUID in URL path');
            return;
        }

        const fetchTrainerGroup = async () => {
            const res = await TrainingsService.getTrainerGroup(trainerUUID, trainingUUID, token)
            setTraining(res);
            console.log(res);
            return res;
        }
        fetchTrainerGroup();
    }, [token, trainerUUID, trainingUUID]);

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


