import React from 'react';

import {Text, Box, Container, Flex} from '@chakra-ui/react';
import './Chat.css';
import {formatDate} from "../../shared/date-time/data-time-pipe";

const ChatHistory = (currentUser, chats) => {
    const history = chats.map(m => {
        // incoming message on left side
        let margin = '0%';
        let bgcolor = 'gray.200';
        let textAlign = 'left';

        if (m.from === currentUser) {
            margin = '20%';
            bgcolor = 'green.200';
            textAlign = 'left';
        }

        const ts = new Date(m.timestamp * 1000);

        return (
            <Box
                key={m.id}
                textAlign={textAlign}
                width={'80%'}
                p={2}
                marginTop={2}
                marginBottom={2}
                marginLeft={margin}
                paddingRight={2}
                bg={bgcolor}
                borderRadius={20}
            >
                <Text fontSize={'sm'}> {m.message} </Text>
                <Flex
                    justify={'flex-end'}
                    marginBottom={3}
                    marginRight={1}
                >
                    <Text as={'sub'} fontSize="xs">
                        {' '}
                        {formatDate(ts)}{' '}
                    </Text>
                </Flex>
            </Box>
        );
    });

    return <Container>{history}</Container>;
};

export default ChatHistory;