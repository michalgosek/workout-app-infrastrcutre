import React, { PropsWithChildren } from "react";

import NoTrainingsAvailable from "./no-trainings-availabe";
import { Table } from "react-bootstrap";
import { TrainingGroup } from "../../services/trainings-service";
import { useGetAllTrainings } from "./hooks";

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
                        <th>Group name</th>
                        <th>Description</th>
                        <th>Trainer</th>
                        <th>Date</th>
                        <th>Limit</th>
                        <th>Participants</th>
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

const Trainings: React.FC = () => {
    const trainings = useGetAllTrainings()
    if (!trainings?.length) return (<NoTrainingsAvailable />);
    return (<TrainingsTable trainings={trainings} />);
};

export default Trainings;
