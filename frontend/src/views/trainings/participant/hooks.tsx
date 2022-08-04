import { useEffect, useState } from 'react';

import { ParticipantGroupReadModel } from 'services/models';
import { TrainingsService } from 'services/trainings-service';
import { useAuth0 } from '@auth0/auth0-react';

const useGetAllParticipantGroups = () => {
    const [trainings, setTrainings] = useState<[] | ParticipantGroupReadModel[]>([]);
    const [trainerUUID, setTrainerUUID] = useState("");
    const { user } = useAuth0();

    useEffect(() => {
        const fetchAllParticipantTraininigs = async () => {
            if (!user || !user.sub) return
            const uuid = user.sub as string
            const trainings = await TrainingsService.getAllParticipantGroups(uuid) ?? [];
            setTrainerUUID(uuid);
            setTrainings(trainings);
        }
        fetchAllParticipantTraininigs();
    }, [user])
    return { trainings, trainerUUID }
}

export default useGetAllParticipantGroups;