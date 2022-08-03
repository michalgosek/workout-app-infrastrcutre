import { TrainerGroup, TrainingGroup, TrainingsService } from '../../services/trainings-service';
import { useEffect, useState } from 'react';

import { useAuth0 } from '@auth0/auth0-react';

const useGetAllTrainerGroups = () => {
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
    }, [])
    return { trainings, trainerUUID }
}

const useGetAllTrainings = () => {
    const [trainings, setTrainings] = useState<TrainingGroup[]>();
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