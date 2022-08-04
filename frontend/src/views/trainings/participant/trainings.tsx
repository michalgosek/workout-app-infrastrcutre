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
    cancelTrainingCallback: (trainerUUID: string, groupUUID: string) => Promise<void>;
};

const ParticipantTrainingsTable: FC<PropsWithChildren<ParticipantTrainingsTableProps>> = ({ trainings, cancelTrainingCallback }) => {
    return (
        <div className={style.rowtext}>
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>training group name</th>
                        <th>description</th>
                        <th>trainer</th>
                        <th>date</th>
                        <th>actions</th>
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
                                    <DropButton onClickHandle={() => cancelTrainingCallback(t.trainer_uuid, t.uuid)} />
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
    const { user } = useAuth0();

    const cancelTrainingCallback = async (trainerUUID: string, groupUUID: string) => {
        if (!user) return;
        const userUUID = user.sub as string;
        const cancelParticipantWorkout = async () => {
            const res = await TrainingsService.cancelParticipantWorkout(userUUID, trainerUUID, groupUUID);
            console.log(res);
        }
        cancelParticipantWorkout();
    };

    return (trainings.length === 0) ? <NoTrainingsAvailable /> : <ParticipantTrainingsTable trainings={trainings} cancelTrainingCallback={cancelTrainingCallback} />;
}

export default ParticipantTrainingGroups;