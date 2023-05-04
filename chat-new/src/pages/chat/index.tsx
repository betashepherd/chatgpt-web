import './chat.css'
import css from '../../App.module.css'
import '../../chatui-theme.css'
import Chat, {Bubble, MessageProps, Progress, toast, useMessages,} from '@chatui/core'
import '@chatui/core/dist/index.css'
import '@chatui/core/es/styles/index.less'
import React, { useEffect, useRef, useState} from 'react'
import clipboardy from 'clipboardy'
import MdEditor from "md-editor-rt"
import "md-editor-rt/lib/style.css"
import sanitizeHtml from 'sanitize-html';
import {send_question} from '../../services/port';
import {EventSourcePolyfill} from 'event-source-polyfill';
import {getCookie} from "../../utils/cookie";

const defaultQuickReplies = [
    {
        name: '清空会话',
        isNew: true,
        isHighlight: true,
    },
    {
        name: '复制会话',
        isNew: false,
        isHighlight: true,
    },
]

const initialMessages = [
    {
        type: 'text',
        content: {
            text: '您好，我是AI助理',
        },
        user: {avatar: '/avatar.png'},
    },
]

let chatContext: any[] = []
let answer = '';

function App() {
    const {messages, appendMsg, updateMsg, setTyping, prependMsgs} = useMessages(initialMessages)
    const [percentage, setPercentage] = useState(0)

    const handleFocus = () => {
        setTimeout(() => {
            window.scrollTo(0, document.body.scrollHeight)

        }, 10)
    }


    // clearQuestion 清空文本特殊字符
    function clearQuestion(requestText: string) {
        requestText = requestText.replace(/\s/g, '')
        const punctuation = ',.;!?，。！？、…'
        const runeRequestText = requestText.split('')
        const lastChar = runeRequestText[runeRequestText.length - 1]
        if (punctuation.indexOf(lastChar) < 0) {
            requestText = requestText + '。'
        }
        return requestText
    }

    // clearQuestion 清空文本换行符号
    function clearReply(reply: string) {
        // TODO 清洗回复特殊字符
        return reply
    }

    function handleSend(type: string, val: string) {
        if (percentage > 0) {
            toast.fail('正在等待上一次回复，请稍后')
            return
        }
        if (type === 'text' && val.trim()) {
            appendMsg({
                type: 'text',
                content: {text: val},
                position: 'left',
                user: {avatar: '/user.png'},
            })

            setTyping(true)
            setPercentage(10)
            onGenCode(val)
        }
    }

    function renderMessageContent(msg: MessageProps) {
        const {_id, type, content} = msg

        switch (type) {
            case 'text':
                let text = content.text
                let msgId = _id.toString()
                let isHtml = sanitizeHtml(text) !== text;
                const richTextRegex = /(<[^>]+>)|(```[^`]*```)/gi;
                const isRichText = richTextRegex.test(text);
                return (
                    <Bubble id={msgId}><MdEditor
                        style={{float: 'left'}}
                        modelValue={text} // 要展示的markdown字符串
                        previewOnly={true} // 只展示预览框部分
                    ></MdEditor></Bubble>
                );

            default:
                return null
        }
    }

    async function handleQuickReplyClick(item: { name: string }) {
        if (item.name === '清空会话') {
            answer = '';
            chatContext.splice(0)
            messages.splice(0)
            prependMsgs(messages)
            toast.success('会话已清空', 1500)
        }
        if (item.name === '复制会话') {
            if (messages.length <= 1) {
                return
            }
            const r = messages
                .slice(1)
                .filter((it) => it.type === 'text')
                .map((it) => it.content.text)
                .join('\n')
            //console.log('messages', messages, r)
            await clipboardy.write(r)
            toast.success('会话已复制, 去粘贴分享', 1500)
        }
    }

    async function onGenCode(question: string) {
        question = clearQuestion(question)
        chatContext.push({
            role: 'user',
            content: question,
        })

        const res = await send_question(chatContext);
        if (res.data.code === 200) {
            const evtSource = new EventSourcePolyfill('/chat/reply', {
                headers: {
                    Authorization: "Bearer " + getCookie("mojolicious") // 请求头携带 token
                }
            });
            let n = 0;
            answer = '';
            evtSource.addEventListener("message", function (event: any) {
                n += 1;
                let djs = JSON.parse(event.data)
                if (djs.data == "--_--xfsdkjfkjsdfjdksjfkdsjfksdjkfjsdkdjf") {
                    res.data.data.messages.push({role: "assistant", content: answer})
                    chatContext = res.data.data.messages
                    setPercentage(0)
                    setTyping(false)
                } else {
                    answer += djs.data
                    if (n > 1) {
                        updateMsg(res.data.data.id, {
                            type: 'text',
                            content: {text: answer},
                            user: {avatar: '/avatar.png'},
                        });
                    } else {
                        appendMsg({
                            _id: res.data.data.id,
                            type: 'text',
                            content: {text: answer},
                            user: {avatar: '/avatar.png'},
                        });
                    }
                }
            });
            evtSource.addEventListener("open", function (event: any) {
                console.log("open");
            });
            evtSource.addEventListener("error", function (event: any) {
                console.log("error");
            });

            evtSource.addEventListener("close", function (event: any) {
                console.log("close");
            });
        } else {
            setPercentage(0)
            setTyping(false)
            return toast.fail('请求出错，' + res.data.errorMsg, 1500)
        }
    }

    return (
        <div className={css.app}>
            <Chat
                navbar={{
                    leftContent: {
                        icon: 'chevron-left',
                        title: 'Back',
                    },
                    rightContent: [
                        {
                            icon: 'apps',
                            title: 'Applications',
                        },
                        {
                            icon: 'ellipsis-h',
                            title: 'More',
                        },
                    ],
                    title: '基于ChatGPT的AI助手',
                }}
                messages={messages}
                renderMessageContent={renderMessageContent}
                quickReplies={defaultQuickReplies}
                onQuickReplyClick={handleQuickReplyClick}
                onSend={handleSend}
                onInputFocus={handleFocus}
            />
            <Progress value={percentage}/>
        </div>
    )
}

export default App
