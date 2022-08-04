import { FC, PropsWithChildren, useState } from 'react';

import DropButton from 'components/buttons/drop-button';
import NoTrainingsAvailable from 'views/trainings/no-trainings-availabe';
import { Table } from 'react-bootstrap';
import { TrainerGroupReadModel } from 'services/models';
import { TrainingsService } from 'services/trainings-service';
import { useGetAllTrainerGroups } from 'views/trainings/hooks';

type TrainingsTableProps = {
    trainings: TrainerGroupReadModel[];
    trainerUUID: string;
    deleteTrainingCallback: (groupUUID: string, trainerUUID: string) => Promise<void>;
};

const TrainingsTable: FC<PropsWithChildren<TrainingsTableProps>> = ({ trainings, trainerUUID, deleteTrainingCallback }) => {
    return (
        <div className='row'>
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Training group name</th>
                        <th>description</th>
                        <th>date</th>
                        <th>limit</th>
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
                                <td>{t.date}</td>
                                <td>{t.limit}</td>
                                <td>
                                    <DropButton onClickHandle={() => deleteTrainingCallback(t.uuid, trainerUUID)} />
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};


// react query 
const TrainerTrainingGroups: FC = () => {
    const [signal, setSignal] = useState(false);
    const deleteTrainerGroupCallback = async (groupUUID: string, trainerUUID: string) => {
        await TrainingsService.deleteTrainerGroup(groupUUID, trainerUUID);
        setSignal(prev => !prev);
    };
    const { trainerUUID, trainings } = useGetAllTrainerGroups(signal)
    return (trainings.length === 0) ? <NoTrainingsAvailable /> : <TrainingsTable trainings={trainings} trainerUUID={trainerUUID} deleteTrainingCallback={deleteTrainerGroupCallback} />;
}

export default TrainerTrainingGroups;
