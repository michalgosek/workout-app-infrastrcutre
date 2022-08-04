import { TrainerGroupReadModel, TrainingGroupReadModel } from 'services/models';
import { useEffect, useState } from 'react';

import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';

const useGetAllTrainerGroups = () => {
    const [trainings, setTrainings] = useState<[] | TrainerGroupReadModel[]>([]);
    const [trainerUUID, setTrainerUUID] = useState("");
    const { user } = useAuth0();

    useEffect(() => {
        const fetchAllTraininigs = async () => {
            if (!user || !user.sub) return
            const uuid = user.sub as string
            const trainings = await TrainingsService.getAllTrainerGroups(uuid) ?? [];
            setTrainerUUID(uuid);
            setTrainings(trainings);
        }
        fetchAllTraininigs();
    }, [user]);
    return { trainings, trainerUUID }
}

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


export {
    useGetAllTrainerGroups,
    useGetAllTrainings
};