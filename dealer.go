package botbase

import "gopkg.in/telegram-bot-api.v4"
import "regexp"
import "log"

type serviceMsg struct {
    stopBot bool
}

type MessageDealer interface {
    init(chan<- tgbotapi.MessageConfig, chan<- serviceMsg)
    accept(tgbotapi.Message)
    run()
    name() string
}

type handlerTrigger struct {
    re *regexp.Regexp
    cmd string
}

func (t *handlerTrigger) canHandle(msg tgbotapi.Message) bool {
    if t.re != nil && t.re.MatchString(msg.Text) {
        log.Printf("Message text '%s' matched regexp '%s'", msg.Text, t.re)
        return true
    }
    if msg.IsCommand() && t.cmd == msg.Command() {
        log.Printf("Message text '%s' matched command '%s'", msg.Text, t.cmd)
        return true
    }
    log.Printf("Message text '%s' doesn't match either command '%s' or regexp '%s'", msg.Text, t.cmd, t.re)
    return false
}

type IncomingMessageHandler interface {
    init(chan<- tgbotapi.MessageConfig, chan<- serviceMsg) handlerTrigger
    handleOne(tgbotapi.Message)
    name() string
}

type IncomingMessageDealer struct {
    MessageDealer
    handler IncomingMessageHandler
    trigger handlerTrigger
    inMsgCh chan tgbotapi.Message
}

func NewIncomingMessageDealer(h IncomingMessageHandler) *IncomingMessageDealer {
    d := &IncomingMessageDealer{handler: h}
    return d
}

func (d *IncomingMessageDealer) init(outMsgCh chan<- tgbotapi.MessageConfig, srvCh chan<- serviceMsg) {
    d.trigger = d.handler.init(outMsgCh, srvCh)
    d.inMsgCh = make(chan tgbotapi.Message, 0)
}

func (d *IncomingMessageDealer) accept(msg tgbotapi.Message) {
    if d.trigger.canHandle(msg) {
        d.inMsgCh<- msg
    }
}

func (d *IncomingMessageDealer) run() {
    go func() {
        for msg := range d.inMsgCh {
            d.handler.handleOne(msg)
        }
    }()
}

func (d *IncomingMessageDealer) name() string {
    return d.handler.name()
}


type BaseHandler struct {
    outMsgCh chan<- tgbotapi.MessageConfig
    srvCh chan<- serviceMsg
}


type BackgroundMessageDealer struct {
    MessageDealer
    BaseHandler
}

func (d *BackgroundMessageDealer) init(outMsgCh chan<- tgbotapi.MessageConfig, srvCh chan<- serviceMsg) {
    d.outMsgCh = outMsgCh
    d.srvCh = srvCh
}

func (d *BackgroundMessageDealer) accept(tgbotapi.Message) {
    // doing nothing
}

// MessageDealer::run to be overwritten by concrete implementations
