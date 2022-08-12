import { Alert, ButtonGroup, Container } from 'react-bootstrap';
import notificationService, { NotificationReadModel } from '../services/notification-service';
import { useEffect, useState } from "react";

import DropButton from 'components/buttons/drop-button';
import style from './notifications.module.scss';
import { useAuth0 } from '@auth0/auth0-react';

const dropAllNotifications = (userUUID: string) => {
    const deleteAll = async () => {
        const res = await notificationService.deleteNotifications(userUUID)
        console.log(res);
        window.location.reload();
    }
    deleteAll();
}

const Notifications = () => {
    const [notifications, setNotifications] = useState<NotificationReadModel[]>();
    const { user } = useAuth0();

    useEffect(() => {
        if (!user || !user.sub) {
            return
        }
        const fetchAllNotifications = async () => {
            const res = await notificationService.getAllNotifications(user.sub as string)
            setNotifications(res)
        }
        fetchAllNotifications()
    }, [user])

    return (
        <>
            <Container>
                {
                    notifications?.map(n => (
                        <Alert key={n.training_uuid} variant='info'>
                            <Alert.Heading>{n.title}</Alert.Heading>
                            <hr />
                            <div>
                                <p><b>Date:</b> {n.date}</p>
                                <p><b>Message:</b> {n.content}</p>
                                <p><b>Trainer:</b> {n.trainer}</p>
                            </div>
                        </Alert>
                    ))
                }
            </Container>
            {
                (notifications?.length !== 0) ?
                    <Container className={style["vertical-center"]}>
                        <DropButton onClickHandle={() => dropAllNotifications(user?.sub as string)} />
                    </Container>
                    : <h1>There is no notifications!  </h1>
            }
        </>

    );
};

export default Notifications; 