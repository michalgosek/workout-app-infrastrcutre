import React, { Fragment, PropsWithChildren } from "react";
import { TrainingGroup, TrainingsService } from "../../services/trainings-service";
import { useEffect, useState } from "react";

import { Table } from "react-bootstrap";

interface TrainingsTableProps {
    trainings: TrainingGroup[];
};

const TrainingsTable: React.FC<PropsWithChildren<TrainingsTableProps>> = ({ trainings }) => {
    return (
        <div className="row">
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Training group name</th>
                        <th>description</th>
                        <th>trainer name</th>
                        <th>date</th>
                        <th>limit</th>
                        <th>participants</th>
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
                                <td>{t.limit}</td>
                                <td>{t.participants}</td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};

const NoTrainingsData: React.FC = () => {
    return (
        <div>
            <h1>There is no trainings available.</h1>
        </div>
    )
};

const Trainings: React.FC = () => {
    const [trainings, setTrainings] = useState<[] | TrainingGroup[]>([]);
    useEffect(() => {
        (async () => {
            const trainings = await TrainingsService.getAllTrainings();
            setTrainings(trainings);
        })();
    }, []);
    return ( trainings.length === 0 ? <NoTrainingsData /> : <TrainingsTable trainings={trainings} /> );
};

export default Trainings;
