import { FC, PropsWithChildren } from 'react';
import { TrainerGroup, TrainingsService } from '../../../services/trainings-service';

import NoTrainingsAvailable from "../no-trainings-availabe";
import { Table } from "react-bootstrap";
import { useGetAllTrainerGroups } from '../hooks';

type TrainingsTableProps = {
    trainings: TrainerGroup[];
    trainerUUID: string;
};

const deleteTrainerGroupCallback = async (groupUUID: string, trainerUUID: string) => {
    await TrainingsService.deleteTrainerGroup(groupUUID, trainerUUID);
    window.location.reload();
};


const TrainingsTable: FC<PropsWithChildren<TrainingsTableProps>> = ({ trainings, trainerUUID }) => {
    return (
        <div className="row">
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
                                    <button type="button" className="btn btn-danger" onClick={() => deleteTrainerGroupCallback(t.uuid, trainerUUID)}>Drop</button>
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};

const TrainerTrainingGroups: FC = () => {
    const { trainerUUID, trainings } = useGetAllTrainerGroups()
    return (trainings.length === 0 ? <NoTrainingsAvailable /> : <TrainingsTable trainings={trainings} trainerUUID={trainerUUID} />);
}

export default TrainerTrainingGroups;
