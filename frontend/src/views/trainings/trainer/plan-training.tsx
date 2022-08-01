import "react-datepicker/dist/react-datepicker.css";

import { Fragment, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";

import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import Button from 'react-bootstrap/Button';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import Form from 'react-bootstrap/Form';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';

type FormValues = {
    group_name: string;
    appointment: Date;
};

const TrainingPlanningForm: React.FC = () => {
    const [value, setValue] = useState<Date | null>(new Date());
    const { handleSubmit, register } = useForm<FormValues>();
    const onHandleChange = (newValue: Date | null) => {
        setValue(newValue);
    };
    const onSubmitHandle: SubmitHandler<FormValues> = (data) => {
        alert(JSON.stringify(data));
    }

    return (
        <Form onSubmit={handleSubmit(onSubmitHandle)} className="form">
            <Form.Group className="mb-3" controlId="formBasicEmail">
                <Form.Label>Training group name:</Form.Label>
                <Form.Control
                    type="text"
                    placeholder="Training group name"
                    {...register("group_name", { required: "This value is obligatory!" })}
                />
                <Form.Text className="text-muted">
                    This will be used as group name description.
                </Form.Text>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formBasicDate">
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                    <Stack spacing={3}>
                        <DateTimePicker
                            {...register('appointment')}
                            label="Training appointment"
                            value={value}
                            onChange={onHandleChange}
                            ampm={false}
                            renderInput={(params) => <TextField {...params} />}
                        />
                    </Stack>
                </LocalizationProvider>
            </Form.Group>

            <Button variant="primary" type="submit" >
                Submit
            </Button>
        </Form >

    );
};


const PlanTraining: React.FC = () => {
    return (
        <Fragment>
            <TrainingPlanningForm />
        </Fragment>
    );
};

export default PlanTraining;




