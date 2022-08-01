import React, { FC, PropsWithChildren } from 'react';
import { TrainerGroup, TrainingsService } from '../../../services/trainings-service';
import { useEffect, useState } from 'react';

import { Table } from "react-bootstrap";
import { useAuth0 } from '@auth0/auth0-react';

type TrainingsTableProps = {
    trainings: TrainerGroup[];
    trainerUUID: string;
};

const deleteTrainerGroupCallback = async (groupUUID: string, trainerUUID: string) => {
    await TrainingsService.deleteTrainerGroup(groupUUID, trainerUUID);;
    window.location.reload();
};

const NoTrainingsData: React.FC = () => {
    return (
        <div>
            <h1>There is no trainings available.</h1>
        </div>
    );
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
    const [trainings, setTrainings] = useState<[] | TrainerGroup[]>([]);
    const [trainerUUID, setTrainerUUID] = useState("");
    const { user } = useAuth0();

    useEffect(() => {
        const fetchAllTraininigs = async () => {
            if (!user || !user.sub) return
            const uuid = user.sub as string

            const trainings = await TrainingsService.getAllTrainerGroups(uuid);
            if (!trainings) return

            setTrainerUUID(uuid);
            setTrainings(trainings);
        }
        fetchAllTraininigs();
    })
    return (trainings.length === 0 ? <NoTrainingsData /> : <TrainingsTable trainings={trainings} trainerUUID={trainerUUID} />);
}

export default TrainerTrainingGroups;
