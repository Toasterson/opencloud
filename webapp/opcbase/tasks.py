#!/usr/bin/env python
from __future__ import absolute_import, unicode_literals

from celery import shared_task

from opcbase.models import Realm


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


def dispatch_realm_add(realm_pk):
    realm_add.delay({
        'id': realm_pk
    })


def dispatch_realm_destroy(realm_pk):
    pass


def dispatch_realm_change(realm_pk):
    pass
