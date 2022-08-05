import React, { FC, PropsWithChildren } from 'react';
import { useEffect, useState } from 'react';

import NoTrainingsAvailable from './no-trainings-availabe';
import { Table } from 'react-bootstrap';
import { TrainingGroupReadModel } from 'services/models';
import TrainingScheduleButton from 'components/buttons/training-schedule-button';
import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';

type TrainingsTableProps = {
    trainings: TrainingGroupReadModel[];
};

type TrainingsTableBodyProps = {
    trainingGroups: TrainingGroupReadModel[];
};

const TrainingsTableBody: FC<PropsWithChildren<TrainingsTableBodyProps>> = ({ trainingGroups }) => {
    const { user } = useAuth0();
    const userUUID = user?.sub ?? "";
    return (
        <>
            {trainingGroups.map((training, index) => (
                <tr key={training.uuid}>
                    <td>{index + 1}</td>
                    <td>{training.name}</td>
                    <td>{training.description}</td>
                    <td>{training.trainer_name}</td>
                    <td>{training.date}</td>
                    <td>{training.limit}</td>
                    <td>{training.participants}</td>
                    <td>
                        {
                            training.trainer_uuid === userUUID ?
                                <p>Training Group Owner</p>
                                :
                                training.limit !== 0 ? <TrainingScheduleButton trainerUUID={training.trainer_uuid} trainingUUID={training.uuid} /> : <b>Full</b>
                        }
                    </td>
                </tr>
            ))}
        </>
    );
};

const TrainingsTable: React.FC<PropsWithChildren<TrainingsTableProps>> = ({ trainings }) => {
    return (
        <div className='row text-center'>
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Group name</th>
                        <th>Description</th>
                        <th>Trainer</th>
                        <th>Date</th>
                        <th>Limit</th>
                        <th>Participants</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
                    <TrainingsTableBody trainingGroups={trainings} />
                </tbody>
            </Table>
        </div>
    );
};


const useGetAllTrainings = () => {
    const [trainings, setTrainings] = useState<TrainingGroupReadModel[]>();
    useEffect(() => {
        const getAllTrainings = async () => {
            const res = await TrainingsService.getAllTrainings()
            setTrainings(res)
        }
        getAllTrainings();
    }, [])
    return trainings
}

const TrainingGroups: React.FC = () => {
    const trainings = useGetAllTrainings()
    if (!trainings?.length) return <NoTrainingsAvailable />;
    return <TrainingsTable trainings={trainings} />;
};

export default TrainingGroups;
