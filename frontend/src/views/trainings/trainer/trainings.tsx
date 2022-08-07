import { FC, PropsWithChildren, useEffect, useState } from 'react';

import DropButton from 'components/buttons/drop-button';
import NoTrainingsAvailable from 'views/trainings/no-trainings-availabe';
import { Table } from 'react-bootstrap';
import TrainerGroupDetailsButton from 'components/buttons/training-group-details-button';
import { TrainerGroupReadModel } from 'services/models';
import { TrainingsService } from 'services/trainings-service';
import style from './trainings.module.scss';
import { useGetTrainingsServiceAccessToken } from 'services/hooks';
import { useParams } from 'react-router-dom';

const deleteTrainerGroupCallback = async (groupUUID: string, trainerUUID: string, token: string) => {
    const res = await TrainingsService.deleteTrainerGroup(groupUUID, trainerUUID, token);
    console.log(res);
    window.location.reload();
};

type TrainingsTableProps = {
    trainings: TrainerGroupReadModel[];
    trainerUUID: string;
};

const TrainingsTable: FC<PropsWithChildren<TrainingsTableProps>> = ({ trainings, trainerUUID }) => {
    const token = useGetTrainingsServiceAccessToken();
    return (
        <div className={style.rowtext}>
            <Table striped bordered hover>
                <thead>
                    <tr className={style['t-headers']}>
                        <th>#</th>
                        <th>Group name</th>
                        <th>Description</th>
                        <th>Date</th>
                        <th>Limit</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        trainings.map((t, index) => (
                            <tr className={style['t-rows']} key={t.uuid}>
                                <td>{index + 1}</td>
                                <td>{t.name}</td>
                                <td>{t.description}</td>
                                <td>{t.date}</td>
                                <td>{t.limit}</td>
                                <td className={style['td-buttons-row']}>
                                    <DropButton onClickHandle={() => deleteTrainerGroupCallback(t.uuid, trainerUUID, token)} />
                                    <TrainerGroupDetailsButton trainerUUID={trainerUUID} trainingUUID={t.uuid} />
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </Table>
        </div>
    );
};


const useGetTrainerGroupsData = () => {
    const { trainerUUID } = useParams();
    const token = useGetTrainingsServiceAccessToken();
    const [trainings, setTrainings] = useState<[] | TrainerGroupReadModel[]>([]);

    useEffect(() => {
        if (!trainerUUID) {
            console.error('missing trainer uuid in trainer groups path');
            return;
        }
        if (!token) {
            console.error('fetching trainings service token failure');
            return;
        }
        const fetchAllTraininigs = async () => {
            const trainings = await TrainingsService.getAllTrainerGroups(trainerUUID, token) ?? [];
            setTrainings(trainings);
        }
        fetchAllTraininigs();
    }, [token, trainerUUID]);


    return { trainings, trainerUUID }
}

const TrainerGroups: FC = () => {
    const { trainings, trainerUUID } = useGetTrainerGroupsData();
    if (!trainings || !trainerUUID) {
        console.error('missing trainer trainings or in trainerUUID in groups path');
        return null;
    }
    return (trainings.length === 0) ? <NoTrainingsAvailable /> : <TrainingsTable trainerUUID={trainerUUID} trainings={trainings} />;

}

export default TrainerGroups;


