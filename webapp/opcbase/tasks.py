#!/usr/bin/env python
from __future__ import absolute_import, unicode_literals

from celery import shared_task

from opcbase.models import Realm
from ansible.playbook.play import Play


@shared_task
def realm_add(msg):
    realm = Realm.objects.get(pk=msg['id'])
    realm.state = 'RE'
    realm.save()
    msg['status'] = 'Done'
    return msg


@shared_task
def realm_change(msg):
    pass


@shared_task
def realm_destroy(msg):
    pass


def instance_add(msg):
    play_source = dict(
        name="Ansible Play",
        hosts='localhost',
        gather_facts='no',
        tasks=[
            dict(action=dict(module='shell', args='ls'), register='shell_out'),
            dict(action=dict(module='debug', args=dict(msg='{{shell_out.stdout}}')))
        ]
    )
    play = Play().load(play_source)



@shared_task
def instance_destroy(msg):
    pass
