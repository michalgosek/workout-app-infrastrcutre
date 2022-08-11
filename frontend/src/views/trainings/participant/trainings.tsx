import { useEffect, useState } from 'react';

import DropButton from 'components/buttons/drop-button';
import { FC } from 'react';
import NoTrainingsAvailable from '../no-trainings-availabe';
import { ParticipantGroupReadModel } from 'services/models';
import { Table } from 'react-bootstrap';
import { TrainingsService } from 'services/trainings-service';
import style from './trainings.module.scss';
import { useGetTrainingsServiceAccessToken } from 'services/hooks';
import { useParams } from 'react-router-dom';

const ParticipantTrainingGroups: FC = () => {
    const token = useGetTrainingsServiceAccessToken();
    const { participantUUID } = useParams();
    const [trainings, setTrainings] = useState<ParticipantGroupReadModel[]>();

    useEffect(() => {
        const fetchAllParticipantTraininigs = async () => {
            if (!participantUUID) {
                console.error('missing participant UUID');
                return;
            }
            const trainings = await TrainingsService.getAllParticipantGroups(participantUUID, token);
            setTrainings(trainings);
        }
        fetchAllParticipantTraininigs();
    }, [participantUUID, token])


    const cancelParticipantTraining = (userUUID: string, trainerUUID: string, groupUUID: string, token: string) => {
        const cancelation = async (userUUID: string, trainerUUID: string, groupUUID: string, token: string) => {
            const res = await TrainingsService.cancelParticipantWorkout(userUUID, trainerUUID, groupUUID, token);
            console.log(res);
            window.location.reload();
        };
        cancelation(userUUID, trainerUUID, groupUUID, token);
    };

    if (!trainings || trainings.length === 0 || !token || !participantUUID) return <NoTrainingsAvailable />

    return (
        <div className={style.rowtext}>
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Group name</th>
                        <th>Description</th>
                        <th>Trainer</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        trainings.map((t, index) => (
                            <tr key={t.uuid}>
                                <td>{index + 1}</td>
                                <td>{t.name}</td>
                                <td>{t.description}</td>
                                <td>{t.trainer_name}</td>
                                <td>{t.date}</td>
                                <td>
                                    {< DropButton onClickHandle={() => cancelParticipantTraining(participantUUID, t.trainer_uuid, t.uuid, token)} />}
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );

}

export default ParticipantTrainingGroups;