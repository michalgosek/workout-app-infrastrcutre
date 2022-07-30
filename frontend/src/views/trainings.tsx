import { TrainigGroup, TrainingsServiceAPI } from "../services/trainings-service-api";
import { useEffect, useState } from "react";

import PageLayout from "./page-layout";
import React from "react";
import { Table } from "react-bootstrap";

const Trainings: React.FC = () => {
    const [trainings, setTrainings] = useState<[] | TrainigGroup[]>([]);

    useEffect(() => {
        (async () => {
            const trainings = await TrainingsServiceAPI.GetAllTrainings();
            setTrainings(trainings);
        })();
    }, []);


    return (
        <PageLayout>
            <div className="row">
                <Table striped bordered hover>
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>training group name</th>
                            <th>description</th>
                            <th>trainer name</th>
                            <th>date</th>
                            <th>limit</th>
                            <th>participants</th>
                        </tr>
                    </thead>
                    <tbody>
                        {React.Children.toArray(trainings.map((t, index) => (
                            <tr>
                                <td>{index + 1}</td>
                                <td>{t.name}</td>
                                <td>{t.description}</td>
                                <td>{t.trainer_name}</td>
                                <td>{t.date}</td>
                                <td>{t.limit}</td>
                                <td>{t.participants}</td>
                            </tr>
                        )))}

                    </tbody>

                </Table>
            </div>
        </PageLayout>
    )
};

export default Trainings;
