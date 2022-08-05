import { FC, PropsWithChildren } from 'react';

import DropButton from 'components/buttons/drop-button';
import NoTrainingsAvailable from '../no-trainings-availabe';
import { ParticipantGroupReadModel } from 'services/models';
import { Table } from 'react-bootstrap';
import { TrainingsService } from 'services/trainings-service';
import style from './trainings.module.scss';
import { useAuth0 } from '@auth0/auth0-react';
import useGetAllParticipantGroups from './hooks';

type ParticipantTrainingsTableProps = {
    trainings: ParticipantGroupReadModel[];
};

const cancelParticipantTraining = async (userUUID: string, trainerUUID: string, groupUUID: string) => {
    TrainingsService.cancelParticipantWorkout(userUUID, trainerUUID, groupUUID)
        .then(res => console.log(res))
        .catch(err => console.error(err))

    window.location.reload();
};

const ParticipantTrainingsTable: FC<PropsWithChildren<ParticipantTrainingsTableProps>> = ({ trainings }) => {
    const { user } = useAuth0();
    if (!user) return null;

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
                                    <DropButton onClickHandle={() => cancelParticipantTraining(user.sub ?? '', t.trainer_uuid, t.uuid)} />
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};


const ParticipantTrainingGroups: FC = () => {
    const { trainings } = useGetAllParticipantGroups();
    return (trainings.length === 0) ? <NoTrainingsAvailable /> : <ParticipantTrainingsTable trainings={trainings} />;
}

export default ParticipantTrainingGroups;