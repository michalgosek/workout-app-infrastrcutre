import { useEffect, useState } from 'react';

import { ParticipantGroupReadModel } from 'services/models';
import { TrainingsService } from 'services/trainings-service';
import { useParams } from 'react-router-dom';

const useGetAllParticipantGroups = (token: string) => {
    const { participantUUID } = useParams();
    const [trainings, setTrainings] = useState<[] | ParticipantGroupReadModel[]>([]);

    useEffect(() => {
        const fetchAllParticipantTraininigs = async () => {
            if (!participantUUID) {
                console.error('missing participant UUID');
                return;
            }

            const trainings = await TrainingsService.getAllParticipantGroups(participantUUID, token) ?? [];
            setTrainings(trainings);
        }
        fetchAllParticipantTraininigs();
    }, [participantUUID, token])
    return trainings;
}

export default useGetAllParticipantGroups;